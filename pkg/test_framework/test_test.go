package framework

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartQuick(t *testing.T) {
	options := *newOptions(DefaultOperatorDelay(500 * time.Millisecond))
	sinks, cout, cerr, operator := newSinks()
	_ = os.Unsetenv("RND")
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	test := kube.NewTest(t).Setup()
	err := test.StartOperator()
	// error because make tests-operator-start is not blocking
	assert.NotNil(t, err)

	ns := test.Namespace
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "tests-cluster-start \$\{RND\}"
echo "tests-cluster-stop \$\{RND\}"
echo "tests-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
echo "tests-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
echo "tests-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
`
	cmp += fmt.Sprintf(`export RND=%s
tests-cluster-stop %s
tests-cluster-start %s
tests-cleanup %s %s 1
tests-cluster-stop %s
$`, rnd, rnd, rnd, rnd, ns, rnd)
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	cmp = fmt.Sprintf("tests-operator-start %s %s 1\n", rnd, test.Namespace)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}

func TestStartSlowNoCleanup(t *testing.T) {
	options := *newOptions(
		DefaultEnvAlways(),
		DefaultPrefix("test-sleep05"),
		DefaultNoCleanup(),
		DefaultOperatorDelay(200 * time.Millisecond),
	)
	sinks, cout, cerr, operator := newSinks()
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	test := kube.NewTest(t).Setup()

	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)

	err = test.StartOperator()
	// operator already started
	assert.NotNil(t, err)

	ns := test.Namespace
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
echo "test-sleep05-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-start %s
test-sleep05-operator-stop %s %s 1
$`, rnd, rnd, rnd, ns)
	// intentionally not calling StopOperator(), test.Close() should call it for us
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
	// stdout output of the operator goes to the operator sink
	cmp = fmt.Sprintf("test-sleep05-operator-start %s %s 1\n", rnd, ns)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}

func TestStartSlowWithCleanup(t *testing.T) {
	options := *newOptions(
		DefaultEnvAlways(),
		DefaultPrefix("test-sleep05"),
		DefaultOperatorDelay(200 * time.Millisecond),
	)
	sinks, cout, cerr, operator := newSinks()
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)
	test := kube.NewTest(t).Setup()

	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)

	ns := test.Namespace
	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-cluster-stop \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
echo "test-sleep05-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-stop %s
test-sleep05-cluster-start %s
test-sleep05-operator-stop %s %s 1
test-sleep05-cleanup %s %s 1
test-sleep05-cluster-stop %s
$`, rnd, rnd, rnd, rnd, ns, rnd, ns, rnd)
	// intentionally not calling StopOperator(), test.Close() should call it for us
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)

	assert.Regexp(t, cmp, cout.String())
	// stdout output of the operator goes to the operator sink
	cmp = fmt.Sprintf("test-sleep05-operator-start %s %s 1\n", rnd, ns)
	assert.Equal(t, cmp, operator.String())
	assert.Empty(t, cerr.String())
}

func TestRunArbitraryTarget(t *testing.T) {
	options := *newOptions(
		DefaultEnvAlways(),
		DefaultPrefix("test-sleep05"),
		DefaultOperatorDelay(200 * time.Millisecond),
	)
	sinks, cout, _, _ := newSinks()
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)

	test := kube.NewTest(t).Setup()
	// this will block long enough to register "operator running"
	err := test.StartOperator()
	assert.NoError(t, err)

	err = test.RunTarget("foo")
	assert.NoError(t, err)

	// try again, detecting an error this time
	err = test.RunTarget("bar")
	assert.NotNil(t, err)

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	ns := test.Namespace
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-cluster-stop \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
echo "test-sleep05-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-stop %s
test-sleep05-cluster-start %s
test-sleep05-foo %s %s 1
test-sleep05-bar %s %s 1.*error
`, rnd, rnd, rnd, rnd, ns, rnd, ns)
	if runtime.GOOS == "linux" {
		// I am very sorry, but there does not seem to be a way to tell the GNU make to keep quiet here
		// However, this doesn't print on Jammy and above (GNU Make 4.3+)
		cmp += "(Makefile:.* failed\n)?"
	}
	cmp += fmt.Sprintf(`test-sleep05-operator-stop %s %s 1
test-sleep05-cleanup %s %s 1
test-sleep05-cluster-stop %s
$`, rnd, ns, rnd, ns, rnd)
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
}

func TestRunArbitraryTargetWithEnv(t *testing.T) {
	options := *newOptions(
		DefaultEnvAlways(),
		DefaultPrefix("test-sleep05"),
		DefaultOperatorDelay(200 * time.Millisecond),
	)
	sinks, cout, _, _ := newSinks()
	kube := startHarness(options, sinks, nil)
	assert.NotNil(t, kube)

	test := kube.NewTest(t).Setup()
	assert.NotNil(t, test)

	// RND is not a reserved env. variable (i.e. it's not set in test.Setup())
	err := test.RunTarget("foo", map[string]string{
		"RND": "123654",
	})
	assert.NoError(t, err)

	// Both TEST_OPERATOR_NS and TEST_COUNT are reserved and will not be overwritten
	// Also, RND=456321 will not be overwritten by RND=9 because ordering
	err = test.RunTarget("foo", map[string]string{
		"RND":              "456321",
		"TEST_OPERATOR_NS": "foo",
		"TEST_COUNT":       "8",
	}, map[string]string{
		"RND":        "9",
		"TEST_COUNT": "9",
		"FOO":        "-fighters",
	})
	assert.NoError(t, err)

	rnd, ok := os.LookupEnv("RND")
	assert.Equal(t, true, ok)
	ns := test.Namespace
	cmp := `^echo "export RND=.*
echo "test-sleep05-cluster-start \$\{RND\}"
echo "test-sleep05-cluster-stop \$\{RND\}"
echo "test-sleep05-operator-start \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
sleep 0\.5s
echo "test-sleep05-operator-stop \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
echo "test-sleep05-cleanup \$\{RND\} \$\{TEST_OPERATOR_NS\} \$\{TEST_COUNT\}"
`
	cmp += fmt.Sprintf(`export RND=%s
test-sleep05-cluster-stop %s
test-sleep05-cluster-start %s
test-sleep05-foo 123654 %s 1
test-sleep05-foo 456321 %s 1-fighters
test-sleep05-cleanup %s %s 1
test-sleep05-cluster-stop %s
`, rnd, rnd, rnd, ns, ns, rnd, ns, rnd)
	test.Close()
	err = kube.Close()
	assert.NoError(t, err)
	assert.Regexp(t, cmp, cout.String())
}
