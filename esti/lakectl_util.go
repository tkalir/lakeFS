//nolint:unused
package esti

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/go-openapi/swag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// lakectl tests utility functions
var update = flag.Bool("update", false, "update golden files with results")

var (
	errValueNotUnique   = errors.New("value not unique in mapping")
	errVariableNotFound = errors.New("variable not found")
)

var (
	reTimestamp       = regexp.MustCompile(`timestamp: \d+`)
	reTime            = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} [-+]\d{4} \w{1,4}`)
	reCommitID        = regexp.MustCompile(`[\d|a-f]{64}`)
	reShortCommitID   = regexp.MustCompile(`[\d|a-f]{16}`)
	reChecksum        = regexp.MustCompile(`([\d|a-f]{32})|(0x[0-9A-F]{15})`)
	reEndpoint        = regexp.MustCompile(`https?://\w+(:\d+)?/api/v\d+/`)
	rePhysicalAddress = regexp.MustCompile(`/data/[0-9a-v]{20}/(?:[0-9a-v]{20}(?:,.+)?)?`)
	reVariable        = regexp.MustCompile(`\$\{([^${}]+)}`)
	rePreSignURL      = regexp.MustCompile(`https://\S+\?\S+`)
	reSecretAccessKey = regexp.MustCompile(`secret_access_key: \S{16,128}`)
	reAccessKeyID     = regexp.MustCompile(`access_key_id: AKIA\S{12,124}`)
)

func lakectlLocation() string {
	return viper.GetString("binaries_dir") + "/lakectl"
}

func LakectlWithParams(accessKeyID, secretAccessKey, endPointURL string) string {
	return LakectlWithParamsWithPosixPerms(accessKeyID, secretAccessKey, endPointURL, false)
}

func LakectlWithParamsWithPosixPerms(accessKeyID, secretAccessKey, endPointURL string, withPosixPerms bool) string {
	lakectlCmdline := "LAKECTL_CREDENTIALS_ACCESS_KEY_ID=" + accessKeyID +
		" LAKECTL_CREDENTIALS_SECRET_ACCESS_KEY=" + secretAccessKey +
		" LAKECTL_SERVER_ENDPOINT_URL=" + endPointURL +
		" LAKECTL_EXPERIMENTAL_LOCAL_POSIX_PERMISSIONS_ENABLED=" + strconv.FormatBool(withPosixPerms) +
		" " + lakectlLocation()

	return lakectlCmdline
}

func Lakectl() string {
	return LakectlWithParams(viper.GetString("access_key_id"), viper.GetString("secret_access_key"), viper.GetString("endpoint_url"))
}

func LakectlWithPosixPerms() string {
	return LakectlWithParamsWithPosixPerms(viper.GetString("access_key_id"), viper.GetString("secret_access_key"), viper.GetString("endpoint_url"), true)
}

func runShellCommand(t *testing.T, command string, isTerminal bool) ([]byte, error) {
	t.Helper()
	t.Logf("Run shell command '%s'", command)
	// Assuming linux. Not sure if this is correct
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Env = append(os.Environ(),
		"LAKECTL_INTERACTIVE="+strconv.FormatBool(isTerminal),
	)
	return cmd.CombinedOutput()
}

// expandVariables receives a string with (possibly) variables in the form of {VAR_NAME}, and
// a var-name -> value mapping. It attempts to replace all occurrences of all variables with
// the corresponding values from the map. If all variables in the string has their mapping
// expandVariables is successful and returns the result string. If at least one variable does
// not have a mapping, expandVariables fails and returns error
func expandVariables(s string, vars map[string]string) (string, error) {
	s = reVariable.ReplaceAllStringFunc(s, func(varReference string) string {
		if val, ok := vars[varReference[2:len(varReference)-1]]; ok {
			return val
		}
		return varReference
	})

	if missingVar := reVariable.FindString(s); missingVar != "" {
		return "", fmt.Errorf("%w, %s", errVariableNotFound, missingVar)
	}

	return s, nil
}

// embedVariables replaces run-specific values from a string with generic, normalized
// variables, that can later be expanded by expandVariables.
// It receives a string that may contain some run-specific data (e.g. repo-name), and
// a mapping of variable names to values. It then replaces all the values found in the original
// string with the corresponding variable name, in the format of {VAR_NAME}. This string can later
// be consumed by expandVariables.
//
// Notes:
// - embedVariables will fail if 2 different variables maps to the same value. While this is a possible
// scenario (e.g. a file named 'master' in 'master' branch) it cannot be processed by embedVariables
// - Values are processed from longest to shortest, and so, if certain var value contains another (e.g.
// VAR1 -> "xy", VAR2 -> "xyz"), the longest option will be considered first. As an example, the string
// "wxyz contains xy which is a prefix of xyz" will be embedded as
// "w{VAR2} contains {VAR1} which is a prefix of {VAR2}")
func embedVariables(s string, vars map[string]string) (string, error) {
	revMap := make(map[string]string)
	vals := make([]string, 0, len(vars)) // collecting all vals, which will be used as keys, in order to control iteration order

	for k, v := range vars {
		if _, exist := revMap[v]; exist {
			return "", fmt.Errorf("%w, %s", errValueNotUnique, v)
		}
		revMap[v] = k
		vals = append(vals, v)
	}

	// Sorting the reversed keys (variable values) by descending length in order to handle longer names first
	// This will diminish replacing partial names that were used to construct longer names
	sort.Slice(vals, func(i, j int) bool {
		return len(vals[i]) > len(vals[j])
	})

	for _, val := range vals {
		s = strings.ReplaceAll(s, val, "${"+revMap[val]+"}")
	}

	return s, nil
}

func sanitize(output string, vars map[string]string) string {
	// The order of execution below is important as certain expression can contain others
	// and so, should be handled first
	s := strings.ReplaceAll(output, "\r\n", "\n")
	if _, ok := vars["DATE"]; !ok {
		s = normalizeProgramTimestamp(s)
	}
	s = normalizeEndpoint(s, vars["LAKEFS_ENDPOINT"])
	s = normalizePreSignURL(s)                       // should be after storage and endpoint to enable non pre-sign url on azure
	s = normalizeRandomObjectKey(s, vars["STORAGE"]) // should be after pre-sign on azure in order not to break the pre-sign url
	s = normalizeCommitID(s)
	s = normalizeChecksum(s)
	s = normalizeShortCommitID(s)
	s = normalizeAccessKeyID(s)
	s = normalizeSecretAccessKey(s)
	return s
}

func RunCmdAndVerifySuccessWithFile(t *testing.T, cmd string, isTerminal bool, goldenFile string, vars map[string]string) {
	t.Helper()
	runCmdAndVerifyWithFile(t, cmd, goldenFile, false, isTerminal, vars)
}

func RunCmdAndVerifyContainsText(t *testing.T, cmd string, isTerminal bool, expectedRaw string, vars map[string]string) {
	t.Helper()
	runCmdAndVerifyContainsText(t, cmd, false, isTerminal, expectedRaw, vars)
}

func RunCmdAndVerifyFailureContainsText(t *testing.T, cmd string, isTerminal bool, expectedRaw string, vars map[string]string) {
	t.Helper()
	runCmdAndVerifyContainsText(t, cmd, true, isTerminal, expectedRaw, vars)
}

func runCmdAndVerifyContainsText(t *testing.T, cmd string, expectFail, isTerminal bool, expectedRaw string, vars map[string]string) {
	t.Helper()
	s := sanitize(expectedRaw, vars)
	expected, err := expandVariables(s, vars)
	require.NoError(t, err, "Variable embed failed - %s", err)
	sanitizedResult := runCmd(t, cmd, expectFail, isTerminal, vars)
	require.Contains(t, sanitizedResult, expected)
}

func RunCmdAndVerifyFailureWithFile(t *testing.T, cmd string, isTerminal bool, goldenFile string, vars map[string]string) {
	t.Helper()
	runCmdAndVerifyWithFile(t, cmd, goldenFile, true, isTerminal, vars)
}

func runCmdAndVerifyWithFile(t *testing.T, cmd, goldenFile string, expectFail, isTerminal bool, vars map[string]string) {
	t.Helper()
	goldenFile = "golden/" + goldenFile + ".golden"

	if *update {
		updateGoldenFile(t, cmd, isTerminal, goldenFile, vars)
	} else {
		content, err := os.ReadFile(goldenFile)
		if err != nil {
			t.Fatal("Failed to read", goldenFile, "-", err)
		}
		expected := sanitize(string(content), vars)
		runCmdAndVerifyResult(t, cmd, expectFail, isTerminal, expected, vars)
	}
}

func updateGoldenFile(t *testing.T, cmd string, isTerminal bool, goldenFile string, vars map[string]string) {
	t.Helper()
	result, _ := runShellCommand(t, cmd, isTerminal)
	s := sanitize(string(result), vars)
	s, err := embedVariables(s, vars)
	require.NoError(t, err, "Variable embed failed")
	err = os.WriteFile(goldenFile, []byte(s), 0o600) //nolint: mnd
	require.NoErrorf(t, err, "Failed to write file %s", goldenFile)
}

func RunCmdAndVerifySuccess(t *testing.T, cmd string, isTerminal bool, expected string, vars map[string]string) {
	t.Helper()
	runCmdAndVerifyResult(t, cmd, false, isTerminal, expected, vars)
}

func RunCmdAndVerifyFailure(t *testing.T, cmd string, isTerminal bool, expected string, vars map[string]string) {
	t.Helper()
	runCmdAndVerifyResult(t, cmd, true, isTerminal, expected, vars)
}

func runCmd(t *testing.T, cmd string, expectFail bool, isTerminal bool, vars map[string]string) string {
	t.Helper()
	result, err := runShellCommand(t, cmd, isTerminal)
	if expectFail {
		require.Errorf(t, err, "Expected error in '%s' command did not occur. Output: %s", cmd, string(result))
	} else {
		require.NoErrorf(t, err, "Failed to run '%s' command - %s", cmd, string(result))
	}
	return sanitize(string(result), vars)
}

func runCmdAndVerifyResult(t *testing.T, cmd string, expectFail bool, isTerminal bool, expected string, vars map[string]string) {
	t.Helper()
	expanded, err := expandVariables(expected, vars)
	if err != nil {
		t.Fatalf("Failed to extract variables for: \"%s\": %s", cmd, err)
	}
	sanitizedResult := runCmd(t, cmd, expectFail, isTerminal, vars)

	require.Equalf(t, expanded, sanitizedResult, "Unexpected output for %s command", cmd)
}

func normalizeProgramTimestamp(output string) string {
	s := reTimestamp.ReplaceAllString(output, "timestamp: <TIMESTAMP>")
	return reTime.ReplaceAllString(s, "<DATE> <TIME> <TZ>")
}

func normalizeRandomObjectKey(output string, objectPrefix string) string {
	objectPrefix = strings.TrimPrefix(objectPrefix, "/")
	for _, match := range rePhysicalAddress.FindAllString(output, -1) {
		output = strings.Replace(output, objectPrefix+match, objectPrefix+"/<OBJECT_KEY>", 1)
	}
	return output
}

func normalizeCommitID(output string) string {
	return reCommitID.ReplaceAllString(output, "<COMMIT_ID>")
}

func normalizeShortCommitID(output string) string {
	return reShortCommitID.ReplaceAllString(output, "<COMMIT_ID_16>")
}

func normalizeChecksum(output string) string {
	return reChecksum.ReplaceAllString(output, "<CHECKSUM>")
}

func normalizeEndpoint(output string, endpoint string) string {
	return reEndpoint.ReplaceAllString(output, endpoint)
}

func normalizePreSignURL(output string) string {
	return rePreSignURL.ReplaceAllString(output, "<PRE_SIGN_URL>")
}

func normalizeAccessKeyID(output string) string {
	return reAccessKeyID.ReplaceAllString(output, "access_key_id: <ACCESS_KEY_ID>")
}

func normalizeSecretAccessKey(output string) string {
	return reSecretAccessKey.ReplaceAllString(output, "secret_access_key: <SECRET_ACCESS_KEY>")
}

func GetCommitter(t testing.TB) string {
	t.Helper()
	userResp, err := client.GetCurrentUserWithResponse(context.Background())
	require.NoError(t, err)
	require.NotNil(t, userResp.JSON200)
	committer := userResp.JSON200.User.Id
	email := swag.StringValue(userResp.JSON200.User.Email)
	if email != "" {
		committer = email
	}
	return committer
}

func GetAuthor(t testing.TB) (string, string) {
	t.Helper()
	userResp, err := client.GetCurrentUserWithResponse(context.Background())
	require.NoError(t, err)
	require.NotNil(t, userResp.JSON200)
	author := userResp.JSON200.User.Id
	email := swag.StringValue(userResp.JSON200.User.Email)
	return author, email
}
