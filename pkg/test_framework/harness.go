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

type Harness struct {
	internalState
	harness.Harness
	Options Options
	Sinks   Sinks
	Scheme  *runtime.Scheme
	Client  client.Client
}

func (h *Harness) Close() error {
	stopCluster(h.Options, h.Sinks)
	return nil
}

func (h *Harness) Run(m *testing.M) int {
	defer h.Close()
	return h.Harness.Run(m)
}

func (h *Harness) NewTest(t htesting.T) *Test {
	test := h.Harness.NewTest(t)
	return &Test{
		Test:            *test,
		operatorRunning: false,
		harness:         h,
		testCount:       atomic.AddUint32(&h.internalState.testCounter, 1),
		envs:            map[string]string{},
	}
}

// Borrowed from github.com/dlespiau/kube-test-harness/blob/master/harness.go
func (h *Harness) OpenManifest(manifest string) (*os.File, error) {
	path := filepath.Join(h.Options.ManifestDirectory, manifest)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type Options struct {
	harness.Options
	Makefile      string
	MakeDir       string
	Prefix        string
	OperatorDelay time.Duration
	EnvAlways     bool
}

// Users can use these to capture the "console" output from the spawned sub-processes rather than
// the default os.Stdout and/or os.Stderr . Capturing the output this way might be useful in tests.
type Sinks struct {
	Stdout    []io.Writer
	Stderr    []io.Writer
	Operator  []io.Writer
}

var Kube *Harness

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

type ParseOptionFn func (a* ParseOptions)

func DefaultMakeDir(makedir string) ParseOptionFn {
	return func(a* ParseOptions) {
		a.MakeDir = makedir
	}
}

func DefaultMakefile(makefile string) ParseOptionFn {
	return func(a* ParseOptions) {
		a.Makefile = makefile
	}
}

func DefaultManifests(manifests string) ParseOptionFn {
	return func(a* ParseOptions) {
		a.Manifests = manifests
	}
}

func DefaultPrefix(prefix string) ParseOptionFn {
	return func(a* ParseOptions) {
		a.Prefix = prefix
	}
}

func DefaultNoCleanup() ParseOptionFn {
	return func(a* ParseOptions) {
		a.NoCleanup = true
	}
}

func DefaultOperatorDelay(opdelay time.Duration) ParseOptionFn {
	return func(a* ParseOptions) {
		a.OperatorDelay = opdelay
	}
}

func DefaultEnvAlways() ParseOptionFn {
	return func(a* ParseOptions) {
		a.EnvAlways = true
	}
}

func OverrideOsArgs(osargs []string) ParseOptionFn {
	return func(a* ParseOptions) {
		a.OsArgs = osargs
	}
}

func OverrideCmdLine(cmdline *flag.FlagSet) ParseOptionFn {
	return func(a* ParseOptions) {
		a.CmdLine = cmdline
	}
}

// We are making use of Functional Options pattern here.
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

	// NOTE: We call "sanitize" functions both here and in Start(). This is to enable
	// the users to create Options by hand, in case if they do not want to use this
	// Parse() function e.g. due to command line options processing here.
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

// We have a fair number of optional parameters here, let's use poor man's default
func Start(options Options, sinks* Sinks, scheme* runtime.Scheme) {
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

// NOTE: this function MUST be idempotent, because it will be called both
// when parsing the parameters and when starting the Kube harness
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
		if !generic.Match([]byte(prefix)){
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

func newClient(scheme* runtime.Scheme) client.Client {
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

func startHarness(options Options, sinks Sinks, scheme* runtime.Scheme) *Harness {
	checkMakefile(options, sinks)
	buildEnv(options, sinks)
	stopCluster(options, sinks)
	startCluster(options, sinks)
	return &Harness{
		internalState: internalState{0},
		Harness: *harness.New(options.Options),
		Options: options,
		Sinks:   sinks,
		Scheme:  scheme,
		Client:  nil,
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
