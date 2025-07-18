package gateway

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/treeverse/lakefs/pkg/api/apiutil"
	"github.com/treeverse/lakefs/pkg/auth"
	"github.com/treeverse/lakefs/pkg/catalog"
	gatewayerrors "github.com/treeverse/lakefs/pkg/gateway/errors"
	"github.com/treeverse/lakefs/pkg/gateway/operations"
	"github.com/treeverse/lakefs/pkg/gateway/path"
	"github.com/treeverse/lakefs/pkg/gateway/sig"
	"github.com/treeverse/lakefs/pkg/graveler"
	"github.com/treeverse/lakefs/pkg/httputil"
	"github.com/treeverse/lakefs/pkg/logging"
	"github.com/treeverse/lakefs/pkg/permissions"
	"github.com/treeverse/lakefs/pkg/stats"
)

func AuthenticationHandler(authService auth.GatewayService, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		user, err := auth.GetUser(ctx)
		if err == nil {
			ctx = logging.AddFields(ctx, logging.Fields{logging.UserFieldKey: user.Username})
			req = req.WithContext(auth.WithUser(ctx, user))
			next.ServeHTTP(w, req)
			return
		}
		o := ctx.Value(ContextKeyOperation).(*operations.Operation)
		authenticator := sig.ChainedAuthenticator(
			sig.NewV4Authenticator(req),
			sig.NewV2SigAuthenticator(req, o.FQDN),
			sig.NewJavaV2SigAuthenticator(req, o.FQDN),
		)
		authContext, err := authenticator.Parse()
		if err != nil {
			o.Log(req).WithError(err).Warn("failed to parse signature")
			_ = o.EncodeError(w, req, err, getAPIErrOrDefault(err, gatewayerrors.ErrAccessDenied))
			return
		}
		accessKeyID := authContext.GetAccessKeyID()
		creds, err := authService.GetCredentials(ctx, accessKeyID)
		logger := o.Log(req)
		if err != nil {
			if !errors.Is(err, auth.ErrNotFound) {
				logger.WithError(err).Warn("error getting access key")
				_ = o.EncodeError(w, req, err, gatewayerrors.ErrInternalError.ToAPIErr())
			} else {
				logger.WithError(err).Warn("could not find access key")
				_ = o.EncodeError(w, req, err, gatewayerrors.ErrAccessDenied.ToAPIErr())
			}
			return
		}
		err = authenticator.Verify(creds)
		if err != nil {
			logger.WithError(err).Warn("error verifying credentials for key")
			_ = o.EncodeError(w, req, err, getAPIErrOrDefault(err, gatewayerrors.ErrAccessDenied))
			return
		}

		user, err = authService.GetUser(ctx, creds.Username)
		if err != nil {
			logger.WithError(err).Warn("could not get user for credentials key")
			_ = o.EncodeError(w, req, err, gatewayerrors.ErrAccessDenied.ToAPIErr())
			return
		}
		ctx = logging.AddFields(ctx, logging.Fields{logging.UserFieldKey: user.Username})
		ctx = auth.WithUser(ctx, user)
		ctx = context.WithValue(ctx, ContextKeyAuthContext, authContext)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func EnrichWithParts(bareDomains []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		parts := ParseRequestParts(req.Host, req.URL.Path, bareDomains)
		ctx = context.WithValue(ctx, ContextKeyRepositoryID, parts.Repository)
		ctx = context.WithValue(ctx, ContextKeyRef, parts.Ref)
		ctx = context.WithValue(ctx, ContextKeyPath, parts.Path)
		ctx = context.WithValue(ctx, ContextKeyMatchedHost, parts.MatchedHost)
		req = req.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func getBareDomain(hostname string, bareDomains []string) string {
	for _, bd := range bareDomains {
		if hostname == stripPort(bd) || strings.HasSuffix(hostname, "."+stripPort(bd)) {
			return bd
		}
	}
	// If no matching bare domain found, assume no gateways.s3.domain_name setting existing,
	//  and we're using path-based routing, with whichever domain our Host header specifies.
	return hostname
}

var trailingPortRegexp = regexp.MustCompile(`:\d+$`)

func stripPort(host string) string {
	return trailingPortRegexp.ReplaceAllString(host, "")
}

func EnrichWithOperation(sc *ServerContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		client := httputil.GetRequestLakeFSClient(req)
		o := &operations.Operation{
			Region:            sc.region,
			FQDN:              getBareDomain(stripPort(req.Host), sc.bareDomains),
			Catalog:           sc.catalog,
			MultipartTracker:  sc.multipartTracker,
			BlockStore:        sc.blockStore,
			Auth:              sc.authService,
			VerifyUnsupported: sc.verifyUnsupported,
			Incr: func(action, userID, repository, ref string) {
				logging.FromContext(ctx).
					WithFields(logging.Fields{
						"action":       action,
						"message_type": "action",
						"repository":   repository,
						"ref":          ref,
						"user_id":      userID,
					}).
					Debug("performing S3 action")
				sc.stats.CollectEvent(stats.Event{
					Class:      "s3_gateway",
					Name:       action,
					Repository: repository,
					Ref:        ref,
					UserID:     userID,
					Client:     client,
				})
			},
			PathProvider: sc.pathProvider,
		}
		next.ServeHTTP(w, req.WithContext(context.WithValue(ctx, ContextKeyOperation, o)))
	})
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		o := ctx.Value(ContextKeyOperation).(*operations.Operation)
		start := time.Now()
		mrw := httputil.NewMetricResponseWriter(w)
		httputil.ConcurrentRequests.WithLabelValues("gateway", string(o.OperationID)).Inc()
		defer httputil.ConcurrentRequests.WithLabelValues("gateway", string(o.OperationID)).Dec()
		next.ServeHTTP(mrw, req)
		requestHistograms.WithLabelValues(string(o.OperationID), strconv.Itoa(mrw.StatusCode)).Observe(time.Since(start).Seconds())
	})
}

func EnrichWithRepositoryOrFallback(c *catalog.Catalog, authService auth.GatewayService, fallbackProxy http.Handler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		repoID := ctx.Value(ContextKeyRepositoryID).(string)
		user, _ := auth.GetUser(ctx)
		username := user.Username
		o := ctx.Value(ContextKeyOperation).(*operations.Operation)
		if repoID == "" {
			// action without a repo
			next.ServeHTTP(w, req)
			return
		}
		repo, err := c.GetRepository(ctx, repoID)
		if errors.Is(err, graveler.ErrNotFound) {
			authResp, authErr := authService.Authorize(ctx, &auth.AuthorizationRequest{
				Username: username,
				RequiredPermissions: permissions.Node{
					Permission: permissions.Permission{Action: permissions.ListRepositoriesAction, Resource: "*"},
				},
			})
			if authErr != nil || authResp.Error != nil || !authResp.Allowed {
				_ = o.EncodeError(w, req, err, gatewayerrors.ErrAccessDenied.ToAPIErr())
				return
			}
			if fallbackProxy != nil {
				fallbackProxy.ServeHTTP(w, req)
				return
			}

			// users often set the gateway endpoint in the clients with /api/v1/ which is the openAPI endpoint.
			// returning a more informative error in such case.
			if strings.HasPrefix(req.RequestURI, apiutil.BaseURL) {
				_ = o.EncodeError(w, req, err, gatewayerrors.ErrNoSuchBucketPossibleAPIEndpoint.ToAPIErr())
				return
			}

			_ = o.EncodeError(w, req, err, gatewayerrors.ErrNoSuchBucket.ToAPIErr())
			return
		}
		if repo == nil {
			_ = o.EncodeError(w, req, err, gatewayerrors.ErrInternalError.ToAPIErr())
			return
		}
		req = req.WithContext(context.WithValue(ctx, ContextKeyRepository, repo))
		next.ServeHTTP(w, req)
	})
}

func OperationLookupHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		o := ctx.Value(ContextKeyOperation).(*operations.Operation)
		repoID := ctx.Value(ContextKeyRepositoryID).(string)
		ref := ctx.Value(ContextKeyRef).(string)
		pth := ctx.Value(ContextKeyPath).(string)

		// based on the operation level, we can determine the operation id
		switch {
		case repoID == "":
			o.OperationID = rootBasedOperationID(req.Method)
		case ref != "" && pth != "":
			o.OperationID = pathBasedOperationID(req.Method)
		case ref == "" && pth == "":
			o.OperationID = repositoryBasedOperationID(req.Method)
		default:
			o.OperationID = operations.OperationIDOperationNotFound
		}

		req = req.WithContext(logging.AddFields(ctx, logging.Fields{"operation_id": o.OperationID}))
		next.ServeHTTP(w, req)
	})
}

// memberFold returns true if 'a' is an equal case-folded to a member of bs.
func memberFold(a string, bs []string) bool {
	for _, b := range bs {
		if strings.EqualFold(a, b) {
			return true
		}
	}
	return false
}

type RequestParts struct {
	Repository  string
	Ref         string
	Path        string
	MatchedHost bool
}

// ParseRequestParts returns the repo id, ref and path according to whether the request is path-style or virtual-host-style.
func ParseRequestParts(host string, urlPath string, bareDomains []string) RequestParts {
	var parts RequestParts
	urlPath = strings.TrimPrefix(urlPath, path.Separator)
	var p []string
	ourHosts := httputil.HostsOnly(bareDomains)
	// we need to check using this order:
	// 1. if exact hosts, path based
	// 2. if suffixes, virtual host
	// 3. none of the above, path based
	if memberFold(httputil.HostOnly(host), ourHosts) {
		// path style: extract repo from first part
		p = strings.SplitN(urlPath, path.Separator, 3) //nolint: mnd
		parts.Repository = p[0]
		if len(p) >= 1 {
			p = p[1:]
		}
		parts.MatchedHost = true
	} else {
		// virtual host style: extract repo from subdomain
		host := httputil.HostOnly(host)
		for _, ourHost := range ourHosts {
			if strings.HasSuffix(host, ourHost) {
				parts.Repository = strings.TrimSuffix(host, "."+ourHost)
				parts.MatchedHost = true
				break
			}
		}
		if parts.MatchedHost {
			p = strings.SplitN(urlPath, path.Separator, 2) //nolint: mnd
		}
	}

	if !parts.MatchedHost {
		// assume path-based for domains we don't explicitly know
		p = strings.SplitN(urlPath, path.Separator, 3) //nolint: mnd
		parts.Repository = p[0]
		if len(p) >= 1 {
			p = p[1:]
		}
	}

	// extract ref and path from remaining parts
	if len(p) > 0 {
		parts.Ref = p[0]
	}
	if len(p) > 1 {
		parts.Path = p[1]
	}
	return parts
}

func rootBasedOperationID(method string) operations.OperationID {
	if method == http.MethodGet {
		return operations.OperationIDListBuckets
	}
	return operations.OperationIDOperationNotFound
}

func pathBasedOperationID(method string) operations.OperationID {
	switch method {
	case http.MethodDelete:
		return operations.OperationIDDeleteObject
	case http.MethodPost:
		return operations.OperationIDPostObject
	case http.MethodGet:
		return operations.OperationIDGetObject
	case http.MethodHead:
		return operations.OperationIDHeadObject
	case http.MethodPut:
		return operations.OperationIDPutObject
	default:
		return operations.OperationIDOperationNotFound
	}
}

func repositoryBasedOperationID(method string) operations.OperationID {
	switch method {
	case http.MethodDelete:
		return operations.OperationIDUnsupportedOperation
	case http.MethodPut:
		return operations.OperationIDPutBucket
	case http.MethodHead:
		return operations.OperationIDHeadBucket
	case http.MethodPost:
		return operations.OperationIDDeleteObjects
	case http.MethodGet:
		return operations.OperationIDListObjects
	default:
		return operations.OperationIDOperationNotFound
	}
}
