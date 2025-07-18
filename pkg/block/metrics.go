package block

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/treeverse/lakefs/pkg/httputil"
)

var concurrentOperations = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "blockstore_concurrent_operations",
		Help: "Number of concurrent blockstore operations",
	},
	[]string{"operation", "blockstore_type"},
)

type MetricsAdapter struct {
	adapter Adapter
}

func NewMetricsAdapter(adapter Adapter) Adapter {
	return &MetricsAdapter{adapter: adapter}
}

func (m *MetricsAdapter) InnerAdapter() Adapter {
	return m.adapter
}

func (m *MetricsAdapter) Put(ctx context.Context, obj ObjectPointer, sizeBytes int64, reader io.Reader, opts PutOpts) (*PutResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("put", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("put", m.adapter.BlockstoreType()).Dec()
	return m.adapter.Put(ctx, obj, sizeBytes, reader, opts)
}

func (m *MetricsAdapter) Get(ctx context.Context, obj ObjectPointer) (io.ReadCloser, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("get", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("get", m.adapter.BlockstoreType()).Dec()
	return m.adapter.Get(ctx, obj)
}

func (m *MetricsAdapter) GetWalker(storageID string, opts WalkerOptions) (Walker, error) {
	return m.adapter.GetWalker(storageID, opts)
}

func (m *MetricsAdapter) GetPreSignedURL(ctx context.Context, obj ObjectPointer, mode PreSignMode, filename string) (string, time.Time, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("get_presigned_url", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("get_presigned_url", m.adapter.BlockstoreType()).Dec()
	return m.adapter.GetPreSignedURL(ctx, obj, mode, filename)
}

func (m *MetricsAdapter) GetPresignUploadPartURL(ctx context.Context, obj ObjectPointer, uploadID string, partNumber int) (string, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("get_presign_upload_part_url", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("get_presign_upload_part_url", m.adapter.BlockstoreType()).Dec()
	return m.adapter.GetPresignUploadPartURL(ctx, obj, uploadID, partNumber)
}

func (m *MetricsAdapter) Exists(ctx context.Context, obj ObjectPointer) (bool, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("exists", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("exists", m.adapter.BlockstoreType()).Dec()
	return m.adapter.Exists(ctx, obj)
}

func (m *MetricsAdapter) GetRange(ctx context.Context, obj ObjectPointer, startPosition int64, endPosition int64) (io.ReadCloser, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("get_range", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("get_range", m.adapter.BlockstoreType()).Dec()
	return m.adapter.GetRange(ctx, obj, startPosition, endPosition)
}

func (m *MetricsAdapter) GetProperties(ctx context.Context, obj ObjectPointer) (Properties, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("get_properties", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("get_properties", m.adapter.BlockstoreType()).Dec()
	return m.adapter.GetProperties(ctx, obj)
}

func (m *MetricsAdapter) Remove(ctx context.Context, obj ObjectPointer) error {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("remove", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("remove", m.adapter.BlockstoreType()).Dec()
	return m.adapter.Remove(ctx, obj)
}

func (m *MetricsAdapter) Copy(ctx context.Context, sourceObj, destinationObj ObjectPointer) error {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("copy", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("copy", m.adapter.BlockstoreType()).Dec()
	return m.adapter.Copy(ctx, sourceObj, destinationObj)
}

func (m *MetricsAdapter) CreateMultiPartUpload(ctx context.Context, obj ObjectPointer, r *http.Request, opts CreateMultiPartUploadOpts) (*CreateMultiPartUploadResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("create_multipart_upload", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("create_multipart_upload", m.adapter.BlockstoreType()).Dec()
	return m.adapter.CreateMultiPartUpload(ctx, obj, r, opts)
}

func (m *MetricsAdapter) UploadPart(ctx context.Context, obj ObjectPointer, sizeBytes int64, reader io.Reader, uploadID string, partNumber int) (*UploadPartResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("upload_part", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("upload_part", m.adapter.BlockstoreType()).Dec()
	return m.adapter.UploadPart(ctx, obj, sizeBytes, reader, uploadID, partNumber)
}

func (m *MetricsAdapter) ListParts(ctx context.Context, obj ObjectPointer, uploadID string, opts ListPartsOpts) (*ListPartsResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("list_parts", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("list_parts", m.adapter.BlockstoreType()).Dec()
	return m.adapter.ListParts(ctx, obj, uploadID, opts)
}

func (m *MetricsAdapter) ListMultipartUploads(ctx context.Context, obj ObjectPointer, opts ListMultipartUploadsOpts) (*ListMultipartUploadsResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("list_multipart_uploads", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("list_multipart_uploads", m.adapter.BlockstoreType()).Dec()
	return m.adapter.ListMultipartUploads(ctx, obj, opts)
}

func (m *MetricsAdapter) UploadCopyPart(ctx context.Context, sourceObj, destinationObj ObjectPointer, uploadID string, partNumber int) (*UploadPartResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("upload_copy_part", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("upload_copy_part", m.adapter.BlockstoreType()).Dec()
	return m.adapter.UploadCopyPart(ctx, sourceObj, destinationObj, uploadID, partNumber)
}

func (m *MetricsAdapter) UploadCopyPartRange(ctx context.Context, sourceObj, destinationObj ObjectPointer, uploadID string, partNumber int, startPosition, endPosition int64) (*UploadPartResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("upload_copy_part_range", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("upload_copy_part_range", m.adapter.BlockstoreType()).Dec()
	return m.adapter.UploadCopyPartRange(ctx, sourceObj, destinationObj, uploadID, partNumber, startPosition, endPosition)
}

func (m *MetricsAdapter) AbortMultiPartUpload(ctx context.Context, obj ObjectPointer, uploadID string) error {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("abort_multipart_upload", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("abort_multipart_upload", m.adapter.BlockstoreType()).Dec()
	return m.adapter.AbortMultiPartUpload(ctx, obj, uploadID)
}

func (m *MetricsAdapter) CompleteMultiPartUpload(ctx context.Context, obj ObjectPointer, uploadID string, multipartList *MultipartUploadCompletion) (*CompleteMultiPartUploadResponse, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("complete_multipart_upload", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("complete_multipart_upload", m.adapter.BlockstoreType()).Dec()
	return m.adapter.CompleteMultiPartUpload(ctx, obj, uploadID, multipartList)
}

func (m *MetricsAdapter) BlockstoreType() string {
	return m.adapter.BlockstoreType()
}

func (m *MetricsAdapter) BlockstoreMetadata(ctx context.Context) (*BlockstoreMetadata, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	return m.adapter.BlockstoreMetadata(ctx)
}

func (m *MetricsAdapter) GetStorageNamespaceInfo(storageID string) *StorageNamespaceInfo {
	return m.adapter.GetStorageNamespaceInfo(storageID)
}

func (m *MetricsAdapter) ResolveNamespace(storageID, storageNamespace, key string, identifierType IdentifierType) (QualifiedKey, error) {
	return m.adapter.ResolveNamespace(storageID, storageNamespace, key, identifierType)
}

func (m *MetricsAdapter) GetRegion(ctx context.Context, storageID, storageNamespace string) (string, error) {
	ctx = httputil.SetClientTrace(ctx, m.adapter.BlockstoreType())
	concurrentOperations.WithLabelValues("get_region", m.adapter.BlockstoreType()).Inc()
	defer concurrentOperations.WithLabelValues("get_region", m.adapter.BlockstoreType()).Dec()
	return m.adapter.GetRegion(ctx, storageID, storageNamespace)
}

func (m *MetricsAdapter) RuntimeStats() map[string]string {
	return m.adapter.RuntimeStats()
}
