package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/openzipkin/zipkin-go"
)

var version = "0.0.1"

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

func paasta() (int, error) {
	zipkinURL, _ := os.LookupEnv("PAASTA_ZIPKIN_URL")
	zr, zt, err := initZipkin(zipkinURL)
	if err != nil {
		return 1, err
	}
	defer zr.Close()

	spanEntry := zt.StartSpan("entrypoint")
	defer spanEntry.Finish()

	spanEntryParent := zipkin.Parent(spanEntry.Context())

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
		spanListCommands := zt.StartSpan("list-subcommands", spanEntryParent)
		defer spanListCommands.Finish()

		var err error
		cmds, err := listPaastaCommands()
		if err != nil {
			spanListCommands.Tag("error", err.Error())
			return 1, fmt.Errorf(
				"generating list of sub-commands: %s", err,
			)
		}

		fullCmdName := fmt.Sprintf("paasta-%s", subcommand)
		if _, ok := cmds[fullCmdName]; ok {
			var err error
			subcommandPath, err = exec.LookPath(fullCmdName)
			if err != nil {
				spanListCommands.Tag("error", err.Error())
				return 1, fmt.Errorf(
					"looking up %s in PATH: %s", fullCmdName, err,
				)
			}
		}
		spanListCommands.Finish()
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

	spanExec := zt.StartSpan("exec-subcommand", spanEntryParent)
	spanExec.Tag("args", strings.Join(args, " "))
	spanExec.Tag("subcommandPath", subcommandPath)
	spanExec.Tag("subcommand", subcommand)

	sc := spanExec.Context()
	env := os.Environ()
	env = append(env, fmt.Sprintf("X_B3_TRACE_ID=%v", sc.TraceID))
	env = append(env, fmt.Sprintf("X_B3_SPAN_ID=%v", sc.ID))
	env = append(env, fmt.Sprintf("X_B3_PARENT_ID=%v", sc.ParentID))
	if sc.Sampled != nil && *sc.Sampled {
		env = append(env, "X_B3_SAMPLED=1")
	}
	if sc.Debug {
		env = append(env, "X_B3_FLAGS=1")
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
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Printf("go-paasta: %v\n", version)
		fmt.Printf("zipkin: %v\n", zipkinReporter)
		fmt.Printf("runtime: %v\n", runtime.Version())
		os.Exit(0)
	}

	exit, err := paasta()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	}
	os.Exit(exit)
}
