package framework

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartQuick(t *testing.T) {
	options := *newOptions()
	options.OperatorStartDelay = 500 * time.Millisecond
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks)
	assert.NotNil(t, kube)
	test := kube.NewTest(t)
	err := test.StartOperator()
	// error because make tests-operator-start is not blocking
	assert.NotNil(t, err)
	test.Close()

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
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	cmp = fmt.Sprintf("tests-operator-start %s\n", rnd)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}

func TestStartSlow(t *testing.T) {
	options := *newOptions()
	options.NoCleanup = true
	options.Prefix = "test-sleep05-"
	options.OperatorStartDelay = 200 * time.Millisecond
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks)
	assert.NotNil(t, kube)
	test := kube.NewTest(t)

	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)
	err = test.StartOperator()
	// operator already started
	assert.NotNil(t, err)
	test.Close()

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-start %s
test-sleep05-operator-stop %s
$`, rnd, rnd, rnd)
	// intentionally not calling StopOperator(), kube.Close() should call it for us
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	// stdout output of the operator goes to the operator sink
	cmp = fmt.Sprintf("test-sleep05-operator-start %s\n", rnd)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}
