package framework

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
	appsv1 "k8s.io/api/apps/v1"
)

type Test struct {
	harness.Test

	operatorRunning bool
	harness         *Harness
	testCount       uint32
	envs            map[string]string
}

func (t *Test) Setup() *Test {
	// this bit of defensive programming is to aid unit testing
	if t.harness.Harness.KubeClient() != nil {
		_ = t.Test.Setup()
	}
	t.envs["TEST_OPERATOR_NS"] = t.Namespace
	t.envs["TEST_COUNT"] = fmt.Sprintf("%d", t.testCount)
	return t
}

func (t *Test) StartOperator() error {
	if t.operatorRunning == true {
		return fmt.Errorf("operator already started")
	}
	err := startOperator(t.harness.Options, t.harness.Sinks, t.envs)
	if err == nil {
		t.operatorRunning = true
	}
	return err
}

func (t *Test) StopOperator() {
	if t.operatorRunning {
		stopOperator(t.harness.Options, t.harness.Sinks, t.envs)
		t.operatorRunning = false
	}
}

func (t *Test) RunTarget(name string, env ... map[string]string) error {
	envs := t.envs
	for _, envsInternal := range env {
		for key, val := range envsInternal {
			// Do not overwrite env. set previously, especially those in Setup() above
			if _, ok := envs[key]; !ok {
				envs[key] = val
			}
		}
	}

	return runTarget(t.harness.Options, t.harness.Sinks, name, envs)
}

func (t *Test) DeleteDeployment(d *appsv1.Deployment, timeout time.Duration) {
	t.Test.DeleteDeployment(d)
	t.Test.WaitForDeploymentDeleted(d, timeout)
}

func (t *Test) Close() {
	// If panicking, let Test.Close() do its thing only and keep the operator running
	defer func () {
		t.Test.Close()
	}()
	if r := recover(); r != nil {
		panic(r)
	} else {
		t.StopOperator()
		cleanup(t.harness.Options, t.harness.Sinks, t.envs)
	}
}

// One-shot channel for single error, safe to send() and close() concurrently
// or many times, but only first operation succeeds (others fail silently)
type chanError struct {
	data    chan error
	closing int32
}

func newChanError() *chanError {
	return &chanError{
		make(chan error, 1),
		0,
	}
}

func (c *chanError) send(err error) {
	if atomic.CompareAndSwapInt32(&c.closing, 0, 1) {
		c.data <- err
		// NOTE: we may send a nil error here, this is supported behaviour
		// NOTE: it is recipient responsibility to call close(c.data)
	}
}

func (c *chanError) close() {
	if atomic.CompareAndSwapInt32(&c.closing, 0, 1) {
		close(c.data)
	}
}

type asynchronousHandler struct {
	delay time.Duration
	result error
}

// The logic is not obvious, so some explanation follows:
// when we start the operator process for testing, it is possible that the process
// will fail right away, because of some early-manifest bug. It might also
// for some reason exit prematurely, without reporting an error.
// To discover when this happens, we will wait for the process to return (possibly
// with an error), and will also start a timer to close the channel for the status
// when operatorStartDelay has elapsed.
// If we have received anything on the channel (before it closed), it means that
// the program completed; otherwise we consider it running.
func(h* asynchronousHandler) Handle(cmd *exec.Cmd, wg *sync.WaitGroup) {
	channel := newChanError()
	go func() {
		wg.Wait()
		err := cmd.Wait()
		// will only succeed to send an error if completed before operatorStartDelay
		channel.send(err)
	}()
	go func() {
		time.Sleep(h.delay)
		// safe no-op if the channel closed earlier
		channel.close()
	}()

	// wait on channel.data will complete when either happens:
	// * channel.send(err), i.e. program completed, possibly with error
	// * channel.close(), i.e. Sleep(operatorStartDelay) elapsed
	if err, ok := <-channel.data; ok {
		if err == nil {
			// This will happen if channel.send(nil) was called above, which
			// indicates that the make target to start the operator has
			// exited prematurely, but with success status. This indicates
			// an unknown error, since we expect "make start operator" to block
			// while the operator is running
			err = fmt.Errorf("operator not running")
		}
		h.result = err
		// NOTE: it is recipient responsibility to call close(c.data)
		close(channel.data)
	}
}

func startOperator(options Options, sinks Sinks, envs map[string]string) error {
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.operatorStart()}
	log.Printf("Starting %v ...", args)
	// let's use sinks.Operator as Stdout for operator output
	handler := asynchronousHandler{options.OperatorDelay, nil}
	if err := start(&handler, sinks.Operator,  nil, args, envs); err != nil {
		return err
	}
	return handler.result
}

func stopOperator(options Options, sinks Sinks, envs map[string]string) {
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.operatorStop()}
	log.Printf("Running %v ...", args)
	// allow for errors here
	_ = run(sinks.Stdout, sinks.Stderr, args, envs)
	log.Print("... done")
}

func cleanup(options Options, sinks Sinks, envs map[string]string) {
	if options.NoCleanup {
		log.Printf("Keeping the test objects")
		return
	}
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.cleanup()}
	log.Printf("Running %v ...", args)
	// allow for errors here
	_ = run(sinks.Stdout, sinks.Stderr, args, envs)
	log.Print("... done")
}

func runTarget(options Options, sinks Sinks, name string, envs map[string]string) error {
	makefile := options.Makefile
	makedir := options.MakeDir
	target := options.Prefix + name
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, target}
	log.Printf("Running %v ...", args)
	// allow for errors here
	err := run(sinks.Stdout, sinks.Stderr, args, envs)
	if err == nil {
		log.Print("... done")
	} else {
		log.Printf("error running target %s: %v", target, err)
	}
	return err
}
