package framework

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"testing"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"
	htesting "github.com/dlespiau/kube-test-harness/testing"
	"github.com/subosito/gotenv"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

type internalState struct {
	testCounter uint32
}

// Harness represents the global state of the test framework.
//
// The single global Harness object will be created by Start function (below) and made available through the global
// Kube object. Fields inside this object are meant for read only and should not be modified by the user.
type Harness struct {
	internalState
	harness.Harness
	// Options provided to Start function
	Options Options
	// Sinks provided to Start function. If "nil" was given, then a default
	// empty instance of Sinks will be stored here.
	Sinks Sinks
	// Scheme provided to Start function. This will be "nil" if Start function
	// received a "nil" parameter (i.e there is no default Scheme)
	Scheme *runtime.Scheme
	// Kubernetes client initialised and authenticated for operations on the
	// test cluster, created inside Start function.
	Client client.Client
}

// Close function will stop the test cluster
//
// The call to this function is deferred inside Run, so users do not need to call it directly.
func (h *Harness) Close() error {
	stopCluster(h.Options, h.Sinks)
	return nil
}

// Run function will apply test options and run the Go test cases, using m.Run()
//
// Users are expected to use the return value from this function as a result code of the Go tests, e.g.:
//
//	os.Exit(framework.Kube.Run(m))
//
// It should be called after Start function, which sets the options and starts the test cluster.
func (h *Harness) Run(m *testing.M) int {
	defer h.Close()
	return h.Harness.Run(m)
}

// NewTest will prepare new test case
//
// Each test case needs a small number of extra data:
//   - test namespace, to be used in the K8s cluster by this test case
//   - sequential test number, to aid parallel test execution
//
// This function takes care to prepare both of these. Note, the user should also the call test.Setup
// function to actually create namespace object in the K8s cluster and populate the environment variables
// specific to this test case (which will be used to pass the above data to the "glue" targets)
func (h *Harness) NewTest(t htesting.T) *Test {
	test := h.Harness.NewTest(t)
	return &Test{
		Test:            *test,
		operatorRunning: false,
		harness:         h,
		testCount:       atomic.AddUint32(&h.testCounter, 1),
		envs:            map[string]string{},
	}
}

// OpenManifest can be used to read files bundled with the acceptance tests
//
// These files are expected to reside inside manifests subdirectory. The actual location of this directory
// can be set from the command line with -k8s.manifests option or set with DefaultManifests functional option.
// If not set it will default to "manifests". This path is relative to where the main_test.go file is.
//
// This function is borrowed from github.com/dlespiau/kube-test-harness/blob/master/harness.go
func (h *Harness) OpenManifest(manifest string) (*os.File, error) {
	path := filepath.Join(h.Options.ManifestDirectory, manifest)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Options for test framework
//
// These options can be set from the command line, in which case they will be applied inside Parse function.
//
// It is recommended that MakeDir and EnvAlways are set to sane defaults programmatically using functional options
// DefaultMakeDir and DefaultEnvAlways respectively.
//
// If desired, test code may ignore the command line (by skipping Parse function) and populate Options
// entirely in code, before the call to Start function.
//
// This also includes harness.Options from github.com/dlespiau/kube-test-harness/
type Options struct {
	harness.Options
	// Makefile is the name of the makefile used for "glue" targets required by the test framework.
	// The default name is "Makefile"
	Makefile string
	// MakeDir is the path to the directory containing the Makefile. It should be set relative to main_test.go
	// file by the DefaultMakeDir functional option; this path will be converted inside Parse to absolute path
	// and then stored here. Note, the conversion to absolute path is idempotent and also performed inside
	// Start, which allows the user to pass a relative path to Start here (e.g. if they are not using Parse)
	MakeDir string
	// Prefix is the prefix of the make targets used for "glue" targets executed by the test framework.
	// It defaults to "test". Parse function will "sanitize" this name if it does not end with either of
	// underscore or minus sign by appending a minus sign. Such sanitized name will be stored here.
	// Note, this operation is idempotent and also performed inside Start.
	Prefix string
	// OperatorDelay is the delay for starting operator during tests, during which the test framework
	// will wait for the operator to fail or exit, so it can fail the test case immediately. If the
	// operator continues to run beyond this interval, it is considered to have started successfully and
	// the framework will stop paying attention to its exit status. The default value is 2s
	OperatorDelay time.Duration
	// EnvAlways is a flag used by the operator to decide what to do if the "env" glue target returned
	// environment variables which collide with the environment variables already set before the tests
	// have started. If this flag is set then the variables from the glue target will take priority.
	// It is recommended that tests use the DefaultEnvAlways functional option to enable this functionality.
	EnvAlways bool
}

// Sinks can be used to capture the "console" output from the spawned sub-processes (e.g. glue targets or
// custom targets) for the purpose of testing. Capturing this output may be useful in tests.
type Sinks struct {
	Stdout   []io.Writer
	Stderr   []io.Writer
	Operator []io.Writer
}

// Kube is the global Harness object, created inside Start function.
var Kube *Harness

// ParseOptions is used by Parse function to enable the test code to set own defaults
//
// We are applying Functional Options pattern, described by Rob Pike at
// https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html (first half only)
// and further documented by Dave Cheney at https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// A short version of this pattern is described with Parse function below.
//
// The purpose of functional options in Parse is to customise the default parameters before command line parsing.
// This helps the users to avoid long "go test" invocations with many "-k8s..." parameters, by enabling the test code
// to provide the most useful defaults instead.
//
// Each of the options below is set by some of the Default... functions below, and then applied inside Parse together
// with command line options to create an Options object which can be passed to Start function.
type ParseOptions struct {
	MakeDir       string
	Makefile      string
	Manifests     string
	Prefix        string
	NoCleanup     bool
	OperatorDelay time.Duration
	EnvAlways     bool
	OsArgs        []string
	CmdLine       *flag.FlagSet
}

type ParseOptionFn func(a *ParseOptions)

// DefaultMakeDir should be called by the test code to set the relative path of the "glue" makefile
//
// For example, assuming that the "glue" targets are defined inside the Makefile residing in the project
// root directory and acceptance tests are in the "acceptance" directory, as demonstrated below:
//
//	| Makefile
//	+ acceptance
//	  | main_test.go
//
// then the TestMain function defined inside main_test.go should pass DefaultMakeDir("..") to Parse function.
// This will ensure that framework will be able to find the Makefile in root project directory to execute the
// "glue" targets. The user can override this value using command line option -k8s.makedir
func DefaultMakeDir(makedir string) ParseOptionFn {
	return func(a *ParseOptions) {
		a.MakeDir = makedir
	}
}

// DefaultMakefile can be used to override the default name of the "glue" makefile
//
// The default name is "Makefile", however if the "glue" targets are defined in a different makefile, the test
// code can set a different default using this function. The user can override this option using command line
// option -k8s.makefile
func DefaultMakefile(makefile string) ParseOptionFn {
	return func(a *ParseOptions) {
		a.Makefile = makefile
	}
}

// DefaultManifests can be called to override the default location of manifest files relative to main_test.go
//
// If not set inside the test code, hardcoded subdirectory "manifests" will be used. The user can override this
// value using command line option -k8s.manifests
func DefaultManifests(manifests string) ParseOptionFn {
	return func(a *ParseOptions) {
		a.Manifests = manifests
	}
}

// DefaultPrefix can be called to override the default prefix of "glue" targets inside Makefile
//
// If not set inside the test code, hardcoded prefix "test" will be used. The user can override this value
// using command line option -k8s.prefix. If the the prefix does not end with either of minus or underscore sign
// then the framework will automatically add minus character at the end. Empty prefix is also valid.
//
// The purpose of prefix is to separate targets used as "glue" by tests and other targets which might be defined
// in the same Makefile. For example, in order to prepare environment variables before test cluster is started
// the test framework will invoke "env" glue target, which together with the default prefix makes "test-env"
// target in the Makefile.
func DefaultPrefix(prefix string) ParseOptionFn {
	return func(a *ParseOptions) {
		a.Prefix = prefix
	}
}

// DefaultNoCleanup can be called to disable cleanup after tests
//
// This function has the same effect as command line option -k8s.no-cleanup=true , however it can be overridden by
// the user setting -k8s.no-cleanup=false . If "no cleanup" option is set, then the test framework will not destroy
// test objects in the test K8s cluster and will not shutdown the cluster. This may lead to excessive resource
// utilization by the test cluster and in turn to transient test failures, so it is recommended that this
// function is not used by the test code.
func DefaultNoCleanup() ParseOptionFn {
	return func(a *ParseOptions) {
		a.NoCleanup = true
	}
}

// DefaultOperatorDelay is how long the framework will wait for the operator being tested to report an error
//
// If not set inside the test code, hardcoded 2 seconds will be used. Larger value means that operator can do
// more work before failing, if we want that failure to be detected by the tests immediately inside the
// StartOperator function; it also means that starting the operator will take longer. This can be overridden
// by the user with -k8s.op-delay command line option.
func DefaultOperatorDelay(opdelay time.Duration) ParseOptionFn {
	return func(a *ParseOptions) {
		a.OperatorDelay = opdelay
	}
}

// DefaultEnvAlways should be used to give the "env" glue target priority to override inherited environment
// variables
//
// Normally if the environment variables returned by the "env" glue target are already set in the inherited
// environment (e.g. in the shell where tests are run) then the framework will ignore values set by the target.
// This behaviour might not be desired, since it makes test results dependent on the environment set by the user
// (e.g. they might have KUBECONFIG already set). When "env always" option is set, then the test framework will
// ignore inherited environment variables when parsing "env" glue target output. It is recommended that test code
// should use DefaultEnvAlways(). This can be overridden by the user with command line option -k8s.env-always
func DefaultEnvAlways() ParseOptionFn {
	return func(a *ParseOptions) {
		a.EnvAlways = true
	}
}

// OverrideOsArgs can be used by the test code to override command line parameters passed to tests
//
// Normally tests should not use this function, unless they really want to ignore command line parameters, but
// still want to call Parse function (rather than prepare Options object explicitly in code). This function
// is mainly used for testing of the Parse function.
func OverrideOsArgs(osargs []string) ParseOptionFn {
	return func(a *ParseOptions) {
		a.OsArgs = osargs
	}
}

// OverrideCmdLine can be used by test code to override the default FlagSet used for tests
//
// Normally tests should not use this function. It is only meant for testing of the Parse function.
func OverrideCmdLine(cmdline *flag.FlagSet) ParseOptionFn {
	return func(a *ParseOptions) {
		a.CmdLine = cmdline
	}
}

// Parse function should be called at the start of test suite to parse the command line options provided by the user
//
// Test code may set the default values of the test options using Default... functional options above, e.g.:
//
//	func TestMain(m *testing.M) {
//		options := framework.Parse(
//			framework.DefaultMakeDir(".."),
//			framework.DefaultEnvAlways(),
//		)
//
// In particular, DefaultMakeDir should be used to point to the Makefile directory where "glue" targets are defined.
//
// This function does not have to be called, e.g. if the test code does not accept command line options.
// In this case MakeDir should be set directly inside the Options object and it will be transformed into
// absolute path inside the Start function.
//
// We are making use of Functional Options pattern here, which works as follows:
//   - Parse function takes a variadic slice of functions matching ParseOptionFn signature
//   - each of these functions is responsible for adjusting a different field of ParseOptions structure
//   - the Default... functional options documented directly above can be used to create the required functions
//   - Parse will execute all these functions, hence adjusting default values of the command line parameters
//   - when command line parameters are parsed inside Parse, such adjusted default values will be applied
func Parse(opts ...ParseOptionFn) *Options {
	args := ParseOptions{
		MakeDir:       "",
		Makefile:      "Makefile",
		Manifests:     "manifests",
		Prefix:        "test",
		NoCleanup:     false,
		OperatorDelay: 2 * time.Second,
		EnvAlways:     false,
		OsArgs:        os.Args[1:],
		CmdLine:       flag.CommandLine,
	}
	for _, opt := range opts {
		opt(&args)
	}

	noCleanup := args.CmdLine.Bool("k8s.no-cleanup", args.NoCleanup, "should test cleanup after themselves")
	verbose := args.CmdLine.Bool("k8s.log.verbose", false, "turn on more verbose logging")
	makefile := args.CmdLine.String("k8s.makefile", args.Makefile, "makefile for glue targets, relative to makedir")
	makedir := args.CmdLine.String("k8s.makedir", args.MakeDir, "directory to makefile, relative to test code")
	prefix := args.CmdLine.String("k8s.prefix", args.Prefix, "prefix for glue targets in makefile")
	manifests := args.CmdLine.String("k8s.manifests", args.Manifests, "directory to K8s manifests, relative to test code")
	delay := args.CmdLine.Duration("k8s.op-delay", args.OperatorDelay, "operator start delay")
	envAlways := args.CmdLine.Bool("k8s.env-always", args.EnvAlways, "should always use environment variables from makefile")
	_ = args.CmdLine.Parse(args.OsArgs)

	// NOTE: We call "sanitize" functions both here and in Start(). This is to enable the users to create Options
	// by hand, in case if they do not want to use this Parse() function e.g. to avoid command line options parsing.
	options := Options{
		Options: harness.Options{
			ManifestDirectory: *manifests,
			NoCleanup:         *noCleanup,
			Logger:            &logger.PrintfLogger{},
		},
		Makefile:      *makefile,
		MakeDir:       sanitizeMakeDir(*makedir),
		Prefix:        sanitizePrefix(*prefix),
		OperatorDelay: *delay,
		EnvAlways:     *envAlways,
	}
	if *verbose {
		options.LogLevel = logger.Debug
	}

	return &options
}

// Start function will prepare and start the test cluster and create the K8s client for operating on it
//
// This is arguably the key function of the test framework, since it performs most work:
//   - validation of all "glue" targets
//   - shutting the previously running cluster (if there was one and "no cleanup" is not set)
//   - starting up a test Kubernetes cluster
//   - creating a controller-runtime/client.Client object for manipulating the test cluster
//   - this client object will be made available for use in test code as framework.Kube.Client
//
// Aside from the regular options, test code may also set:
//   - Sinks, to programmatically receive and handle the output of "glue" targets
//   - runtime.Scheme for CRD used by the test controller-runtime/client.Client object
func Start(options Options, sinks *Sinks, scheme *runtime.Scheme) {
	// NOTE: We call "sanitize" functions both here and in Parse() to avoid
	// strong coupling, i.e. we do not make strong assumption as to the format
	// of MakeDir and Prefix here, hence allowing the user to skip Parse()
	options.MakeDir = sanitizeMakeDir(options.MakeDir)
	options.Prefix = sanitizePrefix(options.Prefix)
	if sinks == nil {
		sinks = &Sinks{}
	}
	Kube = startHarness(options, *sinks, scheme)
	Kube.Client = newClient(scheme)
}

// NOTE: this function MUST be idempotent, because it will be called both when parsing the parameters and when
// starting the Kube harness
func sanitizeMakeDir(makedir string) string {
	if makedir == "" {
		makedir = "."
	}
	result, err := filepath.Abs(makedir)
	if err != nil {
		log.Panic(err)
	}
	return result
}

// NOTE: this function MUST be idempotent (ditto)
func sanitizePrefix(prefix string) string {
	if len(prefix) > 0 {
		// Regexp help is here: https://golang.org/pkg/regexp/syntax/
		generic, _ := regexp.Compile(`^[\w-]+$`)
		trailing, _ := regexp.Compile(`^[\w-]*[-_]$`)
		if !generic.Match([]byte(prefix)) {
			log.Panicf("Invalid k8s.prefix '%s'", prefix)
		} else if !trailing.Match([]byte(prefix)) {
			prefix += "-"
		}
	}
	return prefix
}

// Borrowed from github.com/dlespiau/kube-test-harness/blob/master/harness.go
func newClientConfig(kubeconfig string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

func newClient(scheme *runtime.Scheme) client.Client {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		log.Panicf("KUBECONFIG is empty or not set")
	}

	config, err := newClientConfig(kubeconfig)
	if err != nil {
		log.Panic(err)
	}

	cclient, err := client.New(config, client.Options{
		Scheme: scheme,
		Mapper: nil,
	})
	if err != nil {
		log.Panic(err)
	}

	return cclient
}

func startHarness(options Options, sinks Sinks, scheme *runtime.Scheme) *Harness {
	checkMakefile(options, sinks)
	buildEnv(options, sinks)
	stopCluster(options, sinks)
	startCluster(options, sinks)
	return &Harness{
		internalState: internalState{0},
		Harness:       *harness.New(options.Options),
		Options:       options,
		Sinks:         sinks,
		Scheme:        scheme,
		Client:        nil,
	}
}

func checkMakefile(options Options, sinks Sinks) {
	makefile := options.Makefile
	makedir := options.MakeDir
	check := func(target string) {
		args := []string{"make", "-s", "-f", makefile, "-C", makedir, "--dry-run", target}
		log.Printf("Checking %v ...", args)
		err := run(sinks.Stdout, sinks.Stderr, args, nil)
		if err != nil {
			log.Panicf("error checking target %s: %v", target, err)
		} else {
			log.Print("... done")
		}
	}

	check(options.env())
	check(options.clusterStart())
	if !options.NoCleanup {
		check(options.clusterStop())
	}
	check(options.operatorStart())
	check(options.operatorStop())
	check(options.cleanup())
}

type envScanner struct {
	Out []byte
}

func (d *envScanner) Write(line []byte) (n int, err error) {
	d.Out = append(d.Out, line...)
	return len(line), nil
}

func buildEnv(options Options, sinks Sinks) {
	makefile := options.Makefile
	makedir := options.MakeDir
	exports := envScanner{}
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.env()}
	log.Printf("Running %v ...", args)
	// clone sinks.Stdout and add exports
	cout := append([]io.Writer{}, sinks.Stdout...)
	cout = append(cout, &exports)
	err := run(cout, sinks.Stderr, args, nil)
	if err != nil {
		log.Panic(err)
	}
	log.Print("... done")

	env, err := gotenv.StrictParse(bytes.NewReader(exports.Out))
	if err != nil {
		log.Panic(err)
	}
	for key, val := range env {
		// Empty environment variable looks the same as undefined to
		// the user, so let's treat them the same way here, too
		if old, present := os.LookupEnv(key); options.EnvAlways || !present || old == "" {
			if err := os.Setenv(key, val); err != nil {
				log.Panic(err)
			}
		}
	}
}

func startCluster(options Options, sinks Sinks) {
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.clusterStart()}
	log.Printf("Running %v ...", args)
	err := run(sinks.Stdout, sinks.Stderr, args, nil)
	if err != nil {
		log.Panic(err)
	}
	log.Print("... done")
}

func stopCluster(options Options, sinks Sinks) {
	// Do not stop the cluster if panicking, to enable troubleshooting
	if r := recover(); r != nil {
		log.Printf("Keeping the cluster running for troubleshooting")
		panic(r)
	}
	if options.NoCleanup {
		log.Printf("Keeping the cluster running")
		return
	}
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.clusterStop()}
	log.Printf("Running %v ...", args)
	// if this fails that's perfectly OK - the cluster might not have been running!
	_ = run(sinks.Stdout, sinks.Stderr, args, nil)
	log.Print("... done")
}
