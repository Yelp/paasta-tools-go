package framework

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartQuick(t *testing.T) {
	options := *newOptions()
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, &operator}
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
	assert.Regexp(t, cmp + "$", cout.String())
	cmp = fmt.Sprintf(`^tests-operator-start %s\n`, rnd)
	assert.Regexp(t, cmp + "$", operator.String())
	assert.Empty(t, cerr.String())
}

func TestStartSlow(t *testing.T) {
	options := *newOptions()
	options.Prefix = "test-sleep25-"
	options.NoCleanup = true
	cout := bytes.Buffer{}
	cerr := bytes.Buffer{}
	operator := bytes.Buffer{}
	sinks := Sinks{&cout, &cerr, &operator}
	// NOTE: buildEnv never overwrites existing env. variable
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks)
	assert.NotNil(t, kube)
	test := kube.NewTest(t)

	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)
	test.Close()

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "test-sleep25-cluster-start \$\{RND\}"
echo "test-sleep25-operator-start \$\{RND\}"
sleep 2\.5s
echo "test-sleep25-operator-stop \$\{RND\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep25-cluster-start %s
test-sleep25-operator-stop %s
$`, rnd, rnd, rnd)
	// intentionally not calling StopOperator(), kube.Close() should call it for us
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp + "$", cout.String())
	// stdout output of the operator goes to the operator sink
	cmp = fmt.Sprintf(`^test-sleep25-operator-start %s\n`, rnd)
	assert.Regexp(t, cmp + "$", operator.String())
	assert.Empty(t, cerr.String())
}
