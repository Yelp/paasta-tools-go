package framework

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	harness "github.com/dlespiau/kube-test-harness"
	"github.com/dlespiau/kube-test-harness/logger"
	htesting "github.com/dlespiau/kube-test-harness/testing"
	"github.com/subosito/gotenv"
)

type Harness struct {
	harness.Harness
	Options Options
}

func (h *Harness) Close() error {
	stopCluster(h.Options)
	return nil
}

func (h *Harness) Run(m *testing.M) int {
	defer h.Close()
	return h.Harness.Run(m)
}

func (h *Harness) NewTest(t htesting.T) *Test {
	test := h.Harness.NewTest(t)
	return &Test{
		Test: *test,
		stopOperator: false,
		harness: h,
	}
}

type Options struct {
	harness.Options
	Makefile string
	MakeDir  string
	Prefix   string
}

var Kube *Harness

func Start(m *testing.M) {
	noCleanup := flag.Bool("k8s.no-cleanup", false, "should test cleanup after themselves")
	verbose := flag.Bool("k8s.log.verbose", false, "turn on more verbose logging")
	makefile := flag.String("k8s.makefile", "Makefile", "makefile for cluster manipulation targets, relative to makedir")
	makedir := flag.String("k8s.makedir", "", "directory to makefile")
	prefix := flag.String("k8s.prefix", "test", "prefix for make cluster manipulation targets")
	manifests := flag.String("k8s.manifests", "manifests", "directory to K8s manifests")

	flag.Parse()

	options := Options{
		Options: harness.Options{
			ManifestDirectory: *manifests,
			NoCleanup:         *noCleanup,
			Logger:            &logger.PrintfLogger{},
		},
		Makefile: *makefile,
		MakeDir:  sanitizeMakeDir(*makedir),
		Prefix:   sanitizePrefix(*prefix),
	}
	if *verbose {
		options.LogLevel = logger.Debug
	}

	Kube = startHarness(options)
}

func sanitizeMakeDir(makedir string) string {
	if makedir == "" {
		makedir = "."
	}
	result, err := filepath.Abs(makedir)
	if err != nil {
		log.Panic(err)
		return ""
	}
	return result
}

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

func startHarness(options Options) *Harness {
	checkMakefile(options)
	buildEnv(options)
	stopCluster(options)
	startCluster(options)
	return &Harness{
		Harness: *harness.New(options.Options),
		Options: options,
	}
}

func checkMakefile(options Options) {
	makefile := options.Makefile
	makedir := options.MakeDir
	check := func(target string) {
		args := []string{"make", "-s", "-f", makefile, "-C", makedir, "--dry-run", target}
		log.Printf("Checking %v ...", args)
		err := run(nil, nil, args)
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

func (d *envScanner) Make(dst io.Writer, src io.Reader) outputFn {
	scanner := bufio.NewScanner(src)
	return func() {
		for scanner.Scan() {
			// we need our delimiter back!
			line := append(scanner.Bytes(), '\n')
			d.Out = append(d.Out, line...)
			// don't care if this might fail
			_, _ = dst.Write(line)
		}
	}
}

func buildEnv(options Options) {
	makefile := options.Makefile
	makedir := options.MakeDir
	exports := envScanner{}
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.env()}
	log.Printf("Running %v ...", args)
	err := run(&exports, nil, args)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Print("... done")

	env, err := gotenv.StrictParse(bytes.NewReader(exports.Out))
	if err != nil {
		log.Panic(err)
		return
	}
	for key, val := range env {
		// Empty environment variable looks the same as undefined to
		// the user, so let's treat them the same way here, too
		if old, present := os.LookupEnv(key); !present || old == "" {
			if err := os.Setenv(key, val); err != nil {
				log.Panic(err)
				return
			}
		}
	}
}

func startCluster(options Options) {
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.clusterStart()}
	log.Printf("Running %v ...", args)
	err := run(nil, nil, args)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Print("... done")
}

type pipeDevNull struct{}

func (d *pipeDevNull) Make(dst io.Writer, src io.Reader) outputFn {
	return func() {}
}

func stopCluster(options Options) {
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
	_ = run(&pipeDevNull{}, &pipeDevNull{}, args)
	log.Print("... done")
}
