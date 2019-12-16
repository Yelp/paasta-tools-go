package framework

import (
	"bufio"
	"io"
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
	assert.Equal(t, cwd, sanitizeMakeDir(""))
	assert.Equal(t, cwd, sanitizeMakeDir("."))
	assert.Equal(t, filepath.Join(cwd, "foo"), sanitizeMakeDir("foo"))
	assert.Equal(t, filepath.Join(cwd, ".."), sanitizeMakeDir(".."))
}

func TestSanitizePrefix(t *testing.T) {
	assert.Equal(t, "", sanitizePrefix(""))
	assert.Equal(t, "_", sanitizePrefix("_"))
	assert.Equal(t, "-", sanitizePrefix("-"))
	assert.Equal(t, "abc_", sanitizePrefix("abc_"))
	assert.Equal(t, "abc-", sanitizePrefix("abc-"))
	assert.Equal(t, "abc-", sanitizePrefix("abc"))
	assert.Equal(t, "012_abc-", sanitizePrefix("012_abc"))
	assert.Equal(t, "012-abc-", sanitizePrefix("012-abc"))
	assert.Panics(t, func() {sanitizePrefix(" ")} )
	assert.Panics(t, func() {sanitizePrefix("$")} )
	assert.Panics(t, func() {sanitizePrefix(" abc ")} )
}

type buffer struct{
	Out []byte
}

func (d *buffer) Make(dst io.Writer, src io.Reader) outputFn {
	scanner := bufio.NewScanner(src)
	return func() {
		for scanner.Scan() {
			line := scanner.Bytes()
			d.Out = append(d.Out, line...)
		}
	}
}

// Just run makefile and see we have some output
func TestRunSimple(t *testing.T) {
	cerr := &buffer{}
	cout := &buffer{}
	args := []string{"make", "-s", "-C", "tests", "default"}
	err := run(cout, cerr, args)
	assert.NoError(t, err)
	assert.Equal(t, "default", string(cout.Out))
	assert.Empty(t, string(cerr.Out))
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
	checkMakefile(options)
}

func TestCheckFail(t *testing.T) {
	options := *newOptions()
	// expect fail-close-cluster-stop to fail, not skipped
	options.Prefix = "fail-close-"
	assert.Panics(t, func() { checkMakefile(options) })
	// however, stopCluster() just swallows errors
	stopCluster(options)
}

func TestCheckNoCleanup(t *testing.T) {
	options := *newOptions()
	// expect fail-close-cluster-stop to fail, should be skipped
	options.Prefix = "fail-close-"
	options.NoCleanup = true
	checkMakefile(options)
}

func TestStart(t *testing.T) {
	options := *newOptions()
	kube := startHarness(options)
	assert.NotNil(t, kube)
	value, ok := os.LookupEnv("DUMMY")
	assert.Equal(t, true, ok)
	assert.Equal(t, "tests", value)
	err := kube.Close()
	assert.NoError(t, err)
}

func TestStartNoCleanup(t *testing.T) {
	options := *newOptions()
	options.NoCleanup = true
	// expect fail-close-cluster-stop to fail, should be skipped
	options.Prefix = "fail-close-"
	kube := startHarness(options)
	assert.NotNil(t, kube)
	value, _ := os.LookupEnv("DUMMY")
	assert.Equal(t, "tests", value)
	err := kube.Close()
	assert.NoError(t, err)
}
