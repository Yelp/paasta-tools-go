package framework

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync/atomic"
	"time"

	harness "github.com/dlespiau/kube-test-harness"
)

const (
	// TODO move this to Options
	operatorStartDelay = 2000 * time.Millisecond
)

type Test struct {
	harness.Test

	stopOperator bool
	harness *Harness
}

func (t *Test) StartOperator() error {
	err := startOperator(t.harness.Options, t.harness.Sinks)
	if err == nil {
		t.stopOperator = true
	}
	return err
}

func (t *Test) StopOperator() {
	if t.stopOperator {
		stopOperator(t.harness.Options, t.harness.Sinks)
		t.stopOperator = false
	}
}

func (t *Test) Close() {
	// If panicking, let Test.Close() do its thing only and keep the operator running
	defer t.Test.Close()
	if r := recover(); r != nil {
		panic(r)
	} else {
		t.StopOperator()
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

type asyncOutputCopier struct{}

func (d *asyncOutputCopier) Make(dst io.Writer, src io.Reader) outputFn {
	return func() {
		go func() {
			_, _ = io.Copy(dst, src)
		}()
	}
}

func startOperator(options Options, sinks Sinks) error {
	// The logic is not obvious, so some explanation follows:
	// when we start the operator process for testing, it is possible that the process
	// will fail right away, because of some early-manifest bug.
	// To discover when this happens, we will wait for the process to return an
	// error, and will also start a timer to close the channel for that error
	// when operatorStartDelay has elapsed.
	// If we have received an error on the channel, that means the program completed
	// with error before the channel closed; otherwise we consider it running.
	var result error = nil
	var handler handlerFn = func(cmd *exec.Cmd) {
		channel := newChanError()
		go func() {
			err := cmd.Wait()
			// will only succeed to send an error if completed before operatorStartDelay
			channel.send(err)
		}()
		go func() {
			time.Sleep(operatorStartDelay)
			// safe no-op if the channel closed earlier
			channel.close()
		}()

		// wait on channel.data will complete when either happens:
		// * channel.send(err), i.e. program completed with error
		// * channel.close(), i.e. program completed with success or
		//   Sleep(operatorStartDelay) elapsed
		if err, ok := <-channel.data; ok {
			if err == nil {
				// This will happen if channel.send(nil) was called above, which
				// indicates that the make target to start the operator has
				// exited prematurely, but with success status. This indicates
				// an unknown error, since we expect "make start operator" to block
				// while the operator is running
				err = fmt.Errorf("operator not running")
			}
			result = err
			// NOTE: it is recipient responsibility to call close(c.data)
			close(channel.data)
		}
	}

	makefile := options.Makefile
	makedir := options.MakeDir
	coutFactory := &asyncOutputCopier{}
	cerrFactory := &asyncOutputCopier{}
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.operatorStart()}
	log.Printf("Starting %v ...", args)
	// let's use sinks.Operator as Stdout for operator output
	operatorSinks := sinks
	operatorSinks.Stdout = sinks.Operator
	if err := start(handler, coutFactory, cerrFactory, operatorSinks, args); err != nil {
		return err
	}
	return result
}

func stopOperator(options Options, sinks Sinks) {
	makefile := options.Makefile
	makedir := options.MakeDir
	args := []string{"make", "-s", "-f", makefile, "-C", makedir, options.operatorStop()}
	log.Printf("Running %v ...", args)
	_ = run(&pipeDevNull{}, &pipeDevNull{}, sinks, args)
	log.Print("... done")
}
