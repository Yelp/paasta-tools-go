package framework

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeMakedir(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	assert.NotEmpty(t, cwd)
	assert.Equal(t, cwd, sanitizeMakeDir("."))

	r0 := sanitizeMakeDir("")
	assert.Equal(t, cwd, r0)
	assert.Equal(t, r0, sanitizeMakeDir(r0))

	r1 := sanitizeMakeDir("foo")
	assert.Equal(t, filepath.Join(cwd, "foo"), r1)
	assert.Equal(t, r1, sanitizeMakeDir(r1))

	r2 := sanitizeMakeDir("..")
	assert.Equal(t, filepath.Join(cwd, ".."), r2)
	assert.Equal(t, r2, sanitizeMakeDir(r2))

	r3 := sanitizeMakeDir("/.")
	assert.Equal(t, "/", r3)
	assert.Equal(t, r3, sanitizeMakeDir(r3))
}

func TestSanitizePrefix(t *testing.T) {
	r0 := sanitizePrefix("")
	assert.Equal(t, "", r0)
	assert.Equal(t, r0, sanitizePrefix(r0))

	r1 := sanitizePrefix("_")
	assert.Equal(t, "_", r1)
	assert.Equal(t, r1, sanitizePrefix(r1))

	r2 := sanitizePrefix("-")
	assert.Equal(t, "-", sanitizePrefix("-"))
	assert.Equal(t, r2, sanitizePrefix(r2))

	r3 := sanitizePrefix("abc_")
	assert.Equal(t, "abc_", r3)
	assert.Equal(t, r3, sanitizePrefix(r3))

	r4 := sanitizePrefix("abc-")
	assert.Equal(t, "abc-", r4)
	assert.Equal(t, r4, sanitizePrefix(r4))

	r5 := sanitizePrefix("abc")
	assert.Equal(t, "abc-", r5)
	assert.Equal(t, r5, sanitizePrefix(r5))

	r6 := sanitizePrefix("012_abc")
	assert.Equal(t, "012_abc-", r6)
	assert.Equal(t, r6, sanitizePrefix(r6))

	r7 := sanitizePrefix("012-abc")
	assert.Equal(t, "012-abc-", r7)
	assert.Equal(t, r7, sanitizePrefix(r7))

	assert.Panics(t, func() {sanitizePrefix(" ")} )
	assert.Panics(t, func() {sanitizePrefix("$")} )
	assert.Panics(t, func() {sanitizePrefix(" abc ")} )
}

// Just run makefile with no sinks to capture the output
func TestRunNoOutput(t *testing.T) {
	args := []string{"make", "-s", "-C", "tests", "default"}
	_ = os.Setenv("RND", "BAZ")
	err := run([]io.Writer{}, nil, args)
	assert.NoError(t, err)
	err = run(nil, []io.Writer{}, args)
	assert.NoError(t, err)
}

// Just run makefile and see we have some output
func TestRunSimple(t *testing.T) {
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	args := []string{"make", "-s", "-C", "tests", "default"}
	_ = os.Setenv("RND", "BAZ")
	err := run([]io.Writer{&cout}, []io.Writer{&cerr}, args)
	assert.NoError(t, err)
	assert.Equal(t, "default BAZ\n", cout.String())
	assert.Empty(t, cerr.String())
}

func TestParse(t *testing.T) {
	// Verify default options
	defs := *Parse(OverrideOsArgs([]string{}))
	assert.Equal(t, "manifests", defs.ManifestDirectory)
	assert.Equal(t, false, defs.NoCleanup)
	assert.Equal(t, "Makefile", defs.Makefile)
	assert.Equal(t, sanitizeMakeDir(""), defs.MakeDir)
	assert.Equal(t, "test-", defs.Prefix)
	assert.Equal(t, 2 * time.Second, defs.OperatorDelay)
	assert.Equal(t, false, defs.EnvAlways)

	// Test handling of unknown options
	assert.Panics(t, func() {
		_ = Parse(OverrideOsArgs([]string{"-k8s.no-such-option=0"}))
	})

	// Test individual options (except verbose)
	r1 := defs
	r1.MakeDir = sanitizeMakeDir("foo")
	r1.ManifestDirectory = "baz"
	r1.Prefix = "fizz-"
	r1.OperatorDelay = 5 * time.Second
	r1.NoCleanup = true
	r1.EnvAlways = true

	// Options can be set with Default... functions
	o1 := *Parse(
		DefaultMakeDir("foo"),
		DefaultManifests("baz"),
		DefaultPrefix("fizz"),
		DefaultOperatorDelay(5 * time.Second),
		DefaultNoCleanup(),
		DefaultEnvAlways(),
		OverrideOsArgs([]string{}),
		)

	// Options can be set with command line
	assert.Equal(t, r1, o1)
	o2 := *Parse(
		OverrideOsArgs([]string{
			"-k8s.makedir=foo",
			"-k8s.manifests=baz",
			"-k8s.prefix=fizz",
			"-k8s.op-delay=5s",
			"-k8s.no-cleanup=true",
			"-k8s.env-always=true",
		}))
	assert.Equal(t, r1, o2)

	// Options can be set with Default... functions and overridden from command line
	o3 := *Parse(
		DefaultMakeDir("foo"),
		DefaultManifests("baz"),
		DefaultPrefix("fizz"),
		DefaultOperatorDelay(5 * time.Second),
		DefaultNoCleanup(),
		DefaultEnvAlways(),
		OverrideOsArgs([]string{
			"-k8s.makedir=tests",
			"-k8s.manifests=manifests",
			"-k8s.prefix=tests",
			"-k8s.op-delay=2s",
			"-k8s.no-cleanup=false",
			"-k8s.env-always=false",
		}))
	r2 := defs
	r2.Prefix = "tests-"
	r2.MakeDir = sanitizeMakeDir("tests")
	assert.Equal(t, r2, o3)
}

func newOptions(opts ...ParseOptionFn) *Options {
	// The options in the front are applied first
	opts = append([]ParseOptionFn{
		OverrideOsArgs([]string{}),
		DefaultMakeDir("tests"),
		DefaultPrefix("tests"),
	}, opts ...)
	return Parse(opts...)
}

func TestCheckAll(t *testing.T) {
	options := *newOptions()
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	checkMakefile(options, sinks)
	assert.Regexp(t, `^echo "export RND=.*
echo "tests-cluster-start \$\{RND\}"
echo "tests-cluster-stop \$\{RND\}"
echo "tests-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "tests-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "tests-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
$`, cout.String())
	assert.Empty(t, cerr.String())
	assert.Empty(t, operator.String())
}

func TestCheckFail(t *testing.T) {
	// expect fail-close-cluster-stop to fail, not skipped
	options := *newOptions(DefaultPrefix("fail-close"))
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	assert.Panics(t, func() { checkMakefile(options, sinks) })
	// however, stopCluster() just swallows errors
	stopCluster(options, sinks)
	assert.Regexp(t, `^echo "export RND=.*
echo "fail-close-cluster-start \$\{RND\}"
$`, cout.String())
	assert.NotEmpty(t, cerr.String())
	assert.Empty(t, operator.String())
}

func TestCheckNoCleanup(t *testing.T) {
	// expect fail-close-cluster-stop to fail, should be skipped
	options := *newOptions(DefaultPrefix("fail-close"), DefaultNoCleanup())
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	checkMakefile(options, sinks)
	assert.Regexp(t, `^echo "export RND=.*
echo "fail-close-cluster-start \$\{RND\}"
echo "fail-close-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "fail-close-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "fail-close-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
$`, cout.String())
	assert.Empty(t, cerr.String())
	assert.Empty(t, operator.String())
}

func TestStart(t *testing.T) {
	options := *newOptions()
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "tests-cluster-start \$\{RND\}"
echo "tests-cluster-stop \$\{RND\}"
echo "tests-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "tests-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "tests-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
`
	err := kube.Close()
	assert.NoError(t, err)
	cmp += fmt.Sprintf(`export RND=%s
tests-cluster-stop %s
tests-cluster-start %s
tests-cluster-stop %s
$`, rnd, rnd, rnd, rnd)
	assert.Regexp(t, cmp, cout.String())
	assert.Empty(t, cerr.String())
	assert.Empty(t, operator.String())
}

func TestStartNoCleanup(t *testing.T) {
	// expect fail-close-cluster-stop to fail, should be skipped
	options := *newOptions(DefaultPrefix("fail-close"), DefaultNoCleanup())
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "fail-close-cluster-start \$\{RND\}"
echo "fail-close-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "fail-close-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "fail-close-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
`
	err := kube.Close()
	assert.NoError(t, err)
	cmp += fmt.Sprintf(`export RND=%s
fail-close-cluster-start %s
$`, rnd, rnd)
	assert.Regexp(t, cmp, cout.String())
	assert.Empty(t, cerr.String())
	assert.Empty(t, operator.String())
}
