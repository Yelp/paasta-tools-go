## Acceptance test framework

This test framework is based on and inspired by https://github.com/dlespiau/kube-test-harness

#### Principles of operation

The purpose of this acceptance test framework (for Kubernetes operators) is to provide acceptance tests with a robust
and stable operating environment where operators can be exercised running inside real Kubernetes cluster and where
test code can verify their behaviour programmatically. The choice of actual cluster technology is left to the
user - e.g. built-in tests of the test framework use [k3d](https://github.com/rancher/k3d), while some users are
known to rely on [kind](https://github.com/kubernetes-sigs/kind) instead.
This flexibility is afforded by the concept of "glue" targets, which need to be defined by the user in a makefile
to perform specific actions required by the framework. The specific name and location of the makefile depend on
options provided by the user - makefile and makedir; while the names of the targets depend on prefix which can be
also set by the user. Additionally, special manifests directory can be used to store data required by tests
(e.g. K8s objects definitions), and the location of this directory can be also set by the user.

#### Test execution

The overall execution flow of tests looks like this:
  1. Test code defines `func TestMain(m *testing.M)` inside a `main_test.go` file, as expected by `go test`
  2. Inside `TestMain`, a small number of `framework` function is called in specific order:
     1. `options := framework.Parse(...)` to parse the command line parameters. The ellipsis `...` part should
         be replaced by the default values for parameters created with `framework.Default...` functions (documented below).
     2. `framework.Start(*options, nil, nil)` to start the test K8s cluster and initialise the `framework.Kube.Client`
         object which will be used by test code to access the cluster.
     3. `os.Exit(framework.Kube.Run(m))` to run test cases with `m.Run()` (which will be called indirectly inside `Kube.Run`)
  3. Individual tests are defined inside `*_test.go` files, in functions named like `func Test...(t *testing.T)`, as
     expected by `go test`
     1. each test should use `test := framework.Kube.NewTest(t).Setup()` to create unique test namespace on the K8s
        cluster; this name is made available to test code as `test.Namespace` and is also passed to "glue" targets via
        the environment variable `TEST_OPERATOR_NS`.
     2. tests may start operator for testing with `test.StartOperator()`
     3. operator will be shutdown automatically at the end of test, but users can also stop them explicitly with
        `test.StopOperator()` (and perhaps start again within the same test)
     4. tests may access the cluster with `framework.Kube.Client`, e.g. to read or update state of objects
     5. tests may wait for specific changes to take effect in the cluster with `framework.WaitFor()`
     6. tests may load files from manifests subdirectory with `framework.OpenManifest()` function; these files can
        be used to create K8s objects definitions with the help of `framework.Load...` functions and then created
        in the cluster with `framework.Kube.Client`

A unique namespace defined for each test (as well as test sequential number, see below) can be used to enable parallel
execution of tests, with the regular `go test -parallel=N` syntax. This is subject to all the usual conditions, e.g.:
  * `t.Parallel()` has to be placed in test code
  * "glue" targets need to keep each test separated 

#### Options

The acceptance framework defines `framework.Options` type which is required by the`framework.Start()` function. These
options can be either populated explicitly in test code or created with `framework.Parse(...)` function, which also
parses the command line parameters. In order to set the default values for each command line parameter, test code
can use `framework.Defaults...` functions in a call to `framework.Parse(...)`. Path parameters (i.e. "makedir"
and "manifests") should be set relative to test code directory, i.e. where `main_test.go` resides. They will be
converted internally to absolute paths.

Here is a short list of command line options and corresponding `framework.Default...` to set the default value:

  * `-k8s.makedir`, `DefaultMakeDir` should be set to relative path of the "glue" makefile. If not set, then test
     framework will look for the makefile in the same directory where test code resides.
  * `-k8s.makefile`, `DefaultMakefile` set the default name of the "glue" makefile. If not set, `Makefile` will be used.
  * `-k8s.manifests`, `DefaultManifests` set the relative path to manifest files. If not set `manifests` will be used
  * `-k8s.prefix`, `DefaultPrefix` set the default prefix of "glue" targets inside makefile. This prefix is used to
    determine actual target name; default is `test`. See also explanation below.
  * `-k8s.no-cleanup`, `DefaultNoCleanup` set whether test cluster cleanup should be performed after tests
  * `-k8s.op-delay`, `DefaultOperatorDelay` set how long the framework will wait for operator to start. If operator
     process exists within this time, it will be detected and reported by the `test.StartOperator()` function.
  * `-k8s.env-always`, `DefaultEnvAlways` should be set to give the "env" glue target priority to override
     inherited environment variables.

The purpose of prefix is to separate targets used as "glue" by test framework and other targets which might be defined
in the same makefile. For example, in order to prepare environment variables before start of the test cluster,
the test framework will invoke "env" glue target, which together with the default prefix makes `test-env` target name.
Empty string is a valid input, in which case no prefix will be used.

#### Makefile targets
    
The acceptance tests are expected to provide a `Makefile` which defines how specific actions which are required by
the test framework will be performed. This is achieved by the "glue" targets.

Here is the list of targets, with the default `test` prefix;
  * `test-env`, called inside `framework.Start()`. It is expected to return the list of environment variables which
    will be subsequently set in the test process. The variables should be returned in a format appropriate for
    `/etc/environment` parser and (as a minimum) should contain a `KUBECONFIG` setting for the test Kubernetes cluster.
  * `test-cluster-stop`, called inside `framework.Start()` to stop an old instance of the test cluster
     (if one was running; skipped if "no cleanup" option is set)
  * `test-cluster-start`, called inside `framework.Start()` to create and start a new instance of test cluster
  * `test-operator-start`, called inside `test.StartOperator()` to start the operator process for testing
  * `test-operator-stop`, called either inside `test.StopOperator()` or indirectly in `test.Close()` at the end of
     each test to terminate the operator process
  * `test-cleanup` called inside `test.Close()` to remove the test namespace and all objects in it
  * `test-cluster-stop`, called on successful completion of all tests (skipped if "no cleanup" option is set)
 
All targets have access to the environment variables set initially by "env" target. Targets run by `test` functions 
i.e. "operator start" "operator stop" and "cleanup" (and custom targets, see below) receive additional environment
variables set by `test.Setup()`, which are:
  * `TEST_OPERATOR_NS` which contains the name of the Kubernetes namespace created for this test
  * `TEST_COUNT` which contains the sequential number of the test being run (starting with 1)

Both variables should be used by test targets to ensure that all tests run independently of each other and can be
executed in parallel:
  * test namespace should be passed to the operator, to ensure that it will not operate on any other namespace,
  * sequential test number can be used as a suffix to the operator binary, to ensure that it can be independently
    started and stopped for each test.

The user may also invoke custom targets from makefile with `test.RunTarget(...)` function, e.g. to trigger a
particular error condition on the test cluster etc. This custom target will receive the same set of environment
variables like other targets run by `test` functions and can also receive additional environment variables passed
from the test code (e.g. name of the victim pod). The actual name of the target in the makefile will be determined
by applying the same prefix as other "glue" targets (see explanation above).
