package framework

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"
	htesting "github.com/dlespiau/kube-test-harness/testing"
	"github.com/subosito/gotenv"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

type Harness struct {
	harness.Harness
	Options Options
	Sinks   Sinks
	client  client.Client
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

func (h *Harness) Client() client.Client {
	if h.client == nil {
		log.Panicf("k8s client not initialised")
	}
	return h.client
}

type Options struct {
	harness.Options
	Makefile string
	MakeDir  string
	Prefix   string
	OperatorStartDelay time.Duration
}

// Users can use these to capture the "console" output from the spawned sub-processes rather than
// the default os.Stdout and/or os.Stderr . Capturing the output this way might be useful in tests.
type Sinks struct {
	Stdout    []io.Writer
	Stderr    []io.Writer
	Operator  []io.Writer
}

var Kube *Harness

func Parse() *Options {
	noCleanup := flag.Bool("k8s.no-cleanup", false, "should test cleanup after themselves")
	verbose := flag.Bool("k8s.log.verbose", false, "turn on more verbose logging")
	makefile := flag.String("k8s.makefile", "Makefile", "makefile for cluster manipulation targets, relative to makedir")
	makedir := flag.String("k8s.makedir", "", "directory to makefile")
	prefix := flag.String("k8s.prefix", "test", "prefix for make cluster manipulation targets")
	manifests := flag.String("k8s.manifests", "manifests", "directory to K8s manifests")
	delay := flag.Duration("k8s.op-delay", 2 * time.Second, "operator start delay")

	flag.Parse()

	// NOTE: We call "sanitize" functions both here and in Start(). This is to enable
	// the users to create Options by hand, in case if they do not want to use this
	// Parse() function e.g. due to command line options processing here.
	options := Options{
		Options: harness.Options{
			ManifestDirectory: *manifests,
			NoCleanup:         *noCleanup,
			Logger:            &logger.PrintfLogger{},
		},
		Makefile:           *makefile,
		MakeDir:            sanitizeMakeDir(*makedir),
		Prefix:             sanitizePrefix(*prefix),
		OperatorStartDelay: *delay,
	}
	if *verbose {
		options.LogLevel = logger.Debug
	}

	return &options
}

func Start(m *testing.M, options Options, sinks Sinks) {
	// NOTE: We call "sanitize" functions both here and in Parse() to avoid
	// strong coupling, i.e. we do not make strong assumption as to the format
	// of MakeDir and Prefix here, hence allowing the user to skip Parse()
	options.MakeDir = sanitizeMakeDir(options.MakeDir)
	options.Prefix = sanitizePrefix(options.Prefix)
	Kube = startHarness(options, sinks)
	Kube.client = newClient()
}

func LoadUnstructured(r io.Reader) (*unstructured.Unstructured, error) {
	reader, _, isJson := yaml.GuessJSONStream(r, bytes.MinRead)
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if !isJson {
		tmp, err := yaml.ToJSON(data)
		if err != nil {
			return nil, err
		}
		data = tmp
	}

	result := unstructured.Unstructured{}
	err = result.UnmarshalJSON(data)
	return &result, err
}

func LoadInto(r io.Reader, into interface{}) error {
	if err := yaml.NewYAMLOrJSONDecoder(r, bytes.MinRead).Decode(into); err != nil {
		return err
	}
	return nil
}

type WaitSource func()(runtime.Object, error)

func WaitForNone(timeout time.Duration, from WaitSource) error {
	none := int32(0)
	return WaitFor(&none, timeout, from)
}

func WaitFor(reps* int32, timeout time.Duration, from WaitSource) error {
	wanted := int32(1)
	if reps != nil {
		wanted = *reps
	}

	return wait.Poll(time.Second, timeout, func() (bool, error) {
		current, err := from()
		ready := int32(0)
		if err != nil {
			if !errors.IsNotFound(err) {
				return false, err
			}
			// else let's stick with ready = 0
		} else {
			ready, err = getReady(current)
			if err != nil {
				return false, err
			}
		}
		if ready == wanted {
			return true, nil
		}
		return false, nil
	})
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

func getReady(obj runtime.Object) (int32, error) {
	switch t := obj.(type) {
	case *appsv1.StatefulSet:
		return (*t).Status.ReadyReplicas, nil
	case *appsv1.Deployment:
		return (*t).Status.ReadyReplicas, nil
	case *appsv1.DaemonSet:
		return (*t).Status.NumberReady, nil
	case *appsv1.ReplicaSet:
		return (*t).Status.ReadyReplicas, nil
	case *batchv1.Job:
		return (*t).Status.Active, nil
	case *appsv1.StatefulSetList:
		return int32(len((*t).Items)), nil
	case *appsv1.DeploymentList:
		return int32(len((*t).Items)), nil
	case *appsv1.DaemonSetList:
		return int32(len((*t).Items)), nil
	case *appsv1.ReplicaSetList:
		return int32(len((*t).Items)), nil
	case *appsv1.ControllerRevisionList:
		return int32(len((*t).Items)), nil
	case *batchv1.JobList:
		return int32(len((*t).Items)), nil
	case *corev1.PersistentVolumeList:
		return int32(len((*t).Items)), nil
	case *corev1.PersistentVolumeClaimList:
		return int32(len((*t).Items)), nil
	case *corev1.PodList:
		return int32(len((*t).Items)), nil
	case *corev1.ServiceList:
		return int32(len((*t).Items)), nil
	case *corev1.ServiceAccountList:
		return int32(len((*t).Items)), nil
	case *corev1.EndpointsList:
		return int32(len((*t).Items)), nil
	case *corev1.NodeList:
		return int32(len((*t).Items)), nil
	case *corev1.NamespaceList:
		return int32(len((*t).Items)), nil
	case *corev1.EventList:
		return int32(len((*t).Items)), nil
	case *corev1.SecretList:
		return int32(len((*t).Items)), nil
	case *corev1.ConfigMapList:
		return int32(len((*t).Items)), nil
	case *corev1.ComponentStatusList:
		return int32(len((*t).Items)), nil
	case *rbacv1.RoleBindingList:
		return int32(len((*t).Items)), nil
	case *rbacv1.RoleList:
		return int32(len((*t).Items)), nil
	case *rbacv1.ClusterRoleBindingList:
		return int32(len((*t).Items)), nil
	case *rbacv1.ClusterRoleList:
		return int32(len((*t).Items)), nil
	case *unstructured.UnstructuredList:
		return int32(len((*t).Items)), nil
	case *unstructured.Unstructured:
		if !t.IsList() {
			// Consider single object an equivalent for a list of 1
			return 1, nil
		}
		list, err := t.ToList()
		if err != nil {
			log.Panic(err)
		}
		return int32(len(list.Items)), nil
	default:
		log.Panicf("Unsupported type %v", t.GetObjectKind())
	}
	return 0, nil
}

// Borrowed from github.com/dlespiau/kube-test-harness/blob/master/harness.go
func newClientConfig(kubeconfig string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

func newClient() client.Client {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		log.Panicf("KUBECONFIG is empty or not set")
	}

	config, err := newClientConfig(kubeconfig)
	if err != nil {
		log.Panic(err)
	}

	cclient, err := client.New(config, client.Options{})
	if err != nil {
		log.Panic(err)
	}

	return cclient
}

func startHarness(options Options, sinks Sinks) *Harness {
	checkMakefile(options, sinks)
	buildEnv(options, sinks)
	stopCluster(options, sinks)
	startCluster(options, sinks)
	return &Harness{
		Harness: *harness.New(options.Options),
		Options: options,
		Sinks:   sinks,
		client:  nil,
	}
}

func checkMakefile(options Options, sinks Sinks) {
	makefile := options.Makefile
	makedir := options.MakeDir
	check := func(target string) {
		args := []string{"make", "-s", "-f", makefile, "-C", makedir, "--dry-run", target}
		log.Printf("Checking %v ...", args)
		err := run(sinks.Stdout, sinks.Stderr, args)
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
	err := run(cout, sinks.Stderr, args)
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
		if old, present := os.LookupEnv(key); !present || old == "" {
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
	err := run(sinks.Stdout, sinks.Stderr, args)
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
	_ = run(sinks.Stdout, sinks.Stderr, args)
	log.Print("... done")
}
