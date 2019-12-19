package framework

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"
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

// Just run makefile and see we have some output
func TestRunSimple(t *testing.T) {
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, nil}
	args := []string{"make", "-s", "-C", "tests", "default"}
	_ = os.Setenv("RND", "BAZ")
	err := run( nil, sinks, args)
	assert.NoError(t, err)
	assert.Regexp(t, "^default BAZ\n$", cout.String())
	assert.Empty(t, cerr.String())
}

func newOptions() *Options {
	return &Options{
		Options: harness.Options{
			ManifestDirectory: "",
			NoCleanup:         false,
			Logger:            &logger.PrintfLogger{},
		},
		Makefile: "Makefile",
		MakeDir:  sanitizeMakeDir("tests"),
		Prefix:   sanitizePrefix("tests"),
	}
}

func TestCheckAll(t *testing.T) {
	options := *newOptions()
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, nil}
	checkMakefile(options, sinks)
	assert.Regexp(t, `^echo "export RND=.*
echo "tests-cluster-start \$\{RND\}"
echo "tests-cluster-stop \$\{RND\}"
echo "tests-operator-start \$\{RND\}"
echo "tests-operator-stop \$\{RND\}"
$`, cout.String())
	assert.Empty(t, cerr.String())
}

func TestCheckFail(t *testing.T) {
	options := *newOptions()
	// expect fail-close-cluster-stop to fail, not skipped
	options.Prefix = "fail-close-"
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, nil}
	assert.Panics(t, func() { checkMakefile(options, sinks) })
	// however, stopCluster() just swallows errors
	stopCluster(options, sinks)
	assert.Regexp(t, `^echo "export RND=.*
echo "fail-close-cluster-start \$\{RND\}"
$`, cout.String())
	assert.NotEmpty(t, cerr.String())
}

func TestCheckNoCleanup(t *testing.T) {
	options := *newOptions()
	// expect fail-close-cluster-stop to fail, should be skipped
	options.Prefix = "fail-close-"
	options.NoCleanup = true
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, nil}
	checkMakefile(options, sinks)
	assert.Regexp(t, `^echo "export RND=.*
echo "fail-close-cluster-start \$\{RND\}"
echo "fail-close-operator-start \$\{RND\}"
echo "fail-close-operator-stop \$\{RND\}"
$`, cout.String())
	assert.Empty(t, cerr.String())
}

func TestStart(t *testing.T) {
	options := *newOptions()
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, nil}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks)
	assert.NotNil(t, kube)
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "tests-cluster-start \$\{RND\}"
echo "tests-cluster-stop \$\{RND\}"
echo "tests-operator-start \$\{RND\}"
echo "tests-operator-stop \$\{RND\}"
`
	cmp += fmt.Sprintf(`export RND=%s
tests-cluster-stop %s
tests-cluster-start %s
tests-cluster-stop %s
$`, rnd, rnd, rnd, rnd)
	err := kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp + "$", cout.String())
	assert.Empty(t, cerr.String())
}

func TestStartNoCleanup(t *testing.T) {
	options := *newOptions()
	options.NoCleanup = true
	// expect fail-close-cluster-stop to fail, should be skipped
	options.Prefix = "fail-close-"
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, nil}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks)
	assert.NotNil(t, kube)
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "fail-close-cluster-start \$\{RND\}"
echo "fail-close-operator-start \$\{RND\}"
echo "fail-close-operator-stop \$\{RND\}"
`
	cmp += fmt.Sprintf(`export RND=%s
fail-close-cluster-start %s
$`, rnd, rnd)
	err := kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp + "$", cout.String())
	assert.Empty(t, cerr.String())
}
