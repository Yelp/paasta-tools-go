package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

const endpointURL = "http://zipkin.paasta-pnw-devc.yelp/api/v2/spans"

// map[string]bool is emulating a set
func listPaastaCommands() (map[string]bool, error) {
	cmd := exec.Command("/bin/bash", "-p", "-c", "compgen -A command paasta-")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	cmds := map[string]bool{}
	for scanner.Scan() {
		cmds[scanner.Text()] = true
	}
	return cmds, nil
}

func initZipkin(endpointURL string) (reporter.Reporter, *zipkin.Tracer, error) {
	reporter := reporterhttp.NewReporter(endpointURL)

	localEndpoint, err := zipkin.NewEndpoint("paasta-cli", "localhost:0")
	if err != nil {
		return nil, nil, fmt.Errorf("initializing endpoint: %v", err)
	}

	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing sampler: %v", err)
	}

	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithSampler(sampler),
		zipkin.WithLocalEndpoint(localEndpoint),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing tracer: %v", err)
	}

	return reporter, tracer, err
}

func paasta() (int, error) {
	zr, zt, err := initZipkin(endpointURL)
	if err != nil {
		return 1, err
	}
	defer zr.Close()

	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error initializing zipkin endpoint: %s\n",
			err,
		)
	}

	var subcommand string
	var subcommandPath string
	var args []string

	if len(os.Args) > 1 {
		subcommand = os.Args[1]
		spanListCommands := zt.StartSpan("list-subcommands")
		var err error
		cmds, err := listPaastaCommands()
		if err != nil {
			spanListCommands.Tag("error", err.Error())
			spanListCommands.Finish()
			return 1, fmt.Errorf(
				"generating list of sub-commands: %s", err,
			)
		}
		spanListCommands.Finish()

		spanLookupPath := zt.StartSpan("lookup-subcommand")
		fullCmdName := fmt.Sprintf("paasta-%s", subcommand)
		if _, ok := cmds[fullCmdName]; ok {
			var err error
			subcommandPath, err = exec.LookPath(fullCmdName)
			if err != nil {
				spanListCommands.Tag("error", err.Error())
				spanLookupPath.Finish()
				return 1, fmt.Errorf(
					"looking up %s in PATH: %s", fullCmdName, err,
				)
			}
		}
		spanLookupPath.Finish()
	}

	if subcommandPath != "" {
		args = []string{fmt.Sprintf("paasta-%v", subcommand)}
		if len(os.Args) > 2 {
			args = append(args, os.Args[2:]...)
		}
	} else {
		subcommandPath = "/opt/venvs/paasta-tools/bin/paasta"
		args = []string{"paasta"}
		args = append(args, os.Args[1:]...)
	}

	spanExec := zt.StartSpan("exec-subcommand")
	spanExec.Tag("args", strings.Join(args, " "))
	spanExec.Tag("subcommandPath", subcommandPath)
	spanExec.Tag("subcommand", subcommand)

	sc := spanExec.Context()
	env := os.Environ()
	env = append(env, fmt.Sprintf("ZIPKIN_TRACE_ID=%v", sc.TraceID))
	env = append(env, fmt.Sprintf("ZIPKIN_SPAN_ID=%v", sc.ID))
	env = append(env, fmt.Sprintf("ZIPKIN_PARENT_ID=%v", sc.ParentID))
	if sc.Sampled != nil && *sc.Sampled {
		env = append(env, "ZIPKIN_SAMPLED=1")
	}
	if sc.Debug {
		env = append(env, "ZIPKIN_DEBUG=1")
	}

	cmd := &exec.Cmd{
		Path:   subcommandPath,
		Args:   args,
		Env:    env,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		spanExec.Tag("error", err.Error())
		spanExec.Finish()
		return 1, fmt.Errorf("error running %s: %s", subcommandPath, err)
	}
	spanExec.Finish()

	return 0, nil
}

// os.Exit doesn't work well with defered calls
func main() {
	exit, err := paasta()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
	os.Exit(exit)
}
