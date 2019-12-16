package framework

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

const (
	targetEnv   = "env"
	targetClusterStart = "cluster-start"
	targetClusterStop  = "cluster-stop"
	targetOperatorStart   = "operator-start"
	targetOperatorStop = "operator-stop"
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

// This function is used to handle process output streams
type outputFn func()

// Factory used to create instances of outputFn
type outputFnFactory interface {
	Make(dst io.Writer, src io.Reader) outputFn
}

// This function is used to handle the process after it's been started
// e.g. wait for the result etc.
type handlerFn func(cmd *exec.Cmd)

// General purpose wrapper for "exec.Command().Start()". It can be used to:
// * read Stdout (and perhaps parse it)
// * read Stderr (ditto)
// * wait for result (with the right handler)
// * interrupt & kill the process (ditto)
// No logging occurs inside this function. This function will block until
// all 3 functors have finished, so for truly asynchronous execution you
// may want to spawn goroutines inside each.
func start(handler handlerFn, coutFactory outputFnFactory, cerrFactory outputFnFactory, args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	coutReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cerrReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	cout := coutFactory.Make(os.Stdout, coutReader)
	cerr := cerrFactory.Make(os.Stderr, cerrReader)
	err = cmd.Start()
	if err != nil {
		return err
	}
	// otherwise the process started and is running now

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		cout()
		wg.Done()
	}()
	go func() {
		cerr()
		wg.Done()
	}()
	go func() {
		handler(cmd)
		wg.Done()
	}()

	wg.Wait()
	return nil
}

type outputCopier struct{}

func (d *outputCopier) Make(dst io.Writer, src io.Reader) outputFn {
	return func() {
		_, _ = io.Copy(dst, src)
	}
}

// Wrapper for start() function, more specialized synchronous executor
// similar to exec.Command().Run()
func run(coutFactory outputFnFactory, cerrFactory outputFnFactory, args []string) error {
	if coutFactory == nil {
		coutFactory = &outputCopier{}
	}
	if cerrFactory == nil {
		cerrFactory = &outputCopier{}
	}
	var result error = nil
	var handler handlerFn = func(cmd *exec.Cmd) {
		result = cmd.Wait()
	}
	if err := start(handler, coutFactory, cerrFactory, args); err != nil {
		return err
	}
	return result
}
