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
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	test := kube.NewTest(t).Setup()
	err := test.StartOperator()
	// error because make tests-operator-start is not blocking
	assert.NotNil(t, err)
	ns, nset := os.LookupEnv("TEST_OPERATOR_NS")
	assert.Equal(t,true, nset)
	assert.Equal(t, test.Namespace, ns)

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "tests-cluster-start \$\{RND\}"
echo "tests-cluster-stop \$\{RND\}"
echo "tests-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "tests-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "tests-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
`
	cmp += fmt.Sprintf(`export RND=%s
tests-cluster-stop %s
tests-cluster-start %s
tests-cleanup %s %s
tests-cluster-stop %s
$`, rnd, rnd, rnd, rnd, ns, rnd)
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	cmp = fmt.Sprintf("tests-operator-start %s %s\n", rnd, test.Namespace)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}

func TestStartSlowNoCleanup(t *testing.T) {
	options := *newOptions("test-sleep05", "nocleanup")
	options.OperatorStartDelay = 200 * time.Millisecond
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	test := kube.NewTest(t).Setup()

	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)
	ns, nset := os.LookupEnv("TEST_OPERATOR_NS")
	assert.Equal(t, true, nset)
	assert.Equal(t, test.Namespace, ns)
	err = test.StartOperator()
	// operator already started
	assert.NotNil(t, err)
	ns, nset = os.LookupEnv("TEST_OPERATOR_NS")
	assert.Equal(t, true, nset)
	assert.Equal(t, test.Namespace, ns)

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "test-sleep05-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-start %s
test-sleep05-operator-stop %s %s
$`, rnd, rnd, rnd, ns)
	// intentionally not calling StopOperator(), test.Close() should call it for us
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	// stdout output of the operator goes to the operator sink
	cmp = fmt.Sprintf("test-sleep05-operator-start %s %s\n", rnd, ns)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}

func TestStartSlowWithCleanup(t *testing.T) {
	options := *newOptions("test-sleep05")
	options.OperatorStartDelay = 200 * time.Millisecond
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{[]io.Writer{&cout}, []io.Writer{&cerr}, []io.Writer{&operator}}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	test := kube.NewTest(t).Setup()

	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)
	ns, nset := os.LookupEnv("TEST_OPERATOR_NS")
	assert.Equal(t, true, nset)
	assert.Equal(t, test.Namespace, ns)

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-cluster-stop \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\}"
echo "test-sleep05-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-stop %s
test-sleep05-cluster-start %s
test-sleep05-operator-stop %s %s
test-sleep05-cleanup %s %s
test-sleep05-cluster-stop %s
$`, rnd, rnd, rnd, rnd, ns, rnd, ns, rnd)
	// intentionally not calling StopOperator(), test.Close() should call it for us
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	// stdout output of the operator goes to the operator sink
	cmp = fmt.Sprintf("test-sleep05-operator-start %s %s\n", rnd, ns)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}
