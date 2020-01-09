package framework

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"sync"
)

const (
	targetEnv           = "env"
	targetClusterStart  = "cluster-start"
	targetClusterStop   = "cluster-stop"
	targetOperatorStart = "operator-start"
	targetOperatorStop  = "operator-stop"
)

func (o Options) env() string {
	return o.Prefix + targetEnv
}

func (o Options) clusterStart() string {
	return o.Prefix + targetClusterStart
}

func (o Options) clusterStop() string {
	return o.Prefix + targetClusterStop
}

func (o Options) operatorStart() string {
	return o.Prefix + targetOperatorStart
}

func (o Options) operatorStop() string {
	return o.Prefix + targetOperatorStop
}

// This interface is used to handle the process after it's been started
// e.g. wait for the result etc.
type Handler interface {
	Handle(cmd *exec.Cmd)
}

// General purpose wrapper for "exec.Command().Start()". It can be used to:
// * read Stdout both with outSink and sinks (and perhaps parse it)
// * read Stderr with sinks (ditto)
// * wait for result (with the right handler)
// * interrupt & kill the process (ditto)
// No logging occurs inside this function. This function will block until
// all 3 functors have finished, so for truly asynchronous execution you
// may want to spawn goroutines inside each.
func start(handler Handler, outSinks []io.Writer, errSinks []io.Writer, args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	outScan := bufio.NewScanner(outPipe)

	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	errScan := bufio.NewScanner(errPipe)

	wg1 := sync.WaitGroup{}
	wg1.Add(2)
	go func() {
		wg1.Done()
		for outScan.Scan() {
			// we need our delimiter back!
			line := append(outScan.Bytes(), '\n')
			for _, s := range outSinks {
				s.Write(line)
			}
			os.Stdout.Write(line)
		}
	}()
	go func() {
		wg1.Done()
		for errScan.Scan() {
			// we need our delimiter back!
			line := append(errScan.Bytes(), '\n')
			for _, s := range errSinks {
				s.Write(line)
			}
			os.Stderr.Write(line)
		}
	}()

	// We will lose the output of cmd if it is started before Scan() calls in the
	// go routines above. So let's wait here until they are both scheduled.
	wg1.Wait()
	err = cmd.Start()
	if err != nil {
		return err
	}
	// otherwise the process started and is running now

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		handler.Handle(cmd)
		wg2.Done()
	}()

	wg2.Wait()
	return nil
}

type blockingHandler struct {
	result error
}
func(h *blockingHandler) Handle(cmd *exec.Cmd) {
	h.result = cmd.Wait()
}

// Wrapper for start() function, more specialized synchronous executor
// similar to exec.Command().Run()
func run(outSinks []io.Writer, errSinks []io.Writer, args []string) error {
	handler := blockingHandler{}
	if err := start(&handler, outSinks, errSinks, args); err != nil {
		return err
	}
	return handler.result
}
