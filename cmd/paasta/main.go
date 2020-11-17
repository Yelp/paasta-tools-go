package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Yelp/paasta-tools-go/pkg/configstore"
	paastaversion "github.com/Yelp/paasta-tools-go/pkg/version"
	paastazipkin "github.com/Yelp/paasta-tools-go/pkg/zipkin"
	"github.com/openzipkin/zipkin-go"
	"k8s.io/klog"
)

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
	if zipkinURL == "" {
		store := configstore.NewStore(
			"/etc/paasta",
			map[string]string{"paasta_zipkin_url": "paasta"},
		)
		store.Load("paasta_zipkin_url", &zipkinURL)
	}

	zr, zt, err := paastazipkin.InitZipkin(zipkinURL)
	if err != nil {
		klog.V(10).Infof("Error initializing zipkin: %s\n", err)
		err = nil
	}
	defer zr.Close()

	spanEntry := zt.StartSpan("entrypoint")
	defer spanEntry.Finish()

	spanEntryParent := zipkin.Parent(spanEntry.Context())
	if err != nil {
		klog.V(10).Infof("Error initializing zipkin endpoint: %s\n", err)
		err = nil
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
	defer spanExec.Finish()

	spanExec.Tag("args", strings.Join(args, " "))
	spanExec.Tag("subcommandPath", subcommandPath)
	spanExec.Tag("subcommand", subcommand)
	user, ok := os.LookupEnv("SUDO_USER")
	if !ok {
		user, _ = os.LookupEnv("USER")
	}
	spanExec.Tag("user", user)

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
	if subcommand == "-V" {
		fmt.Printf("paasta-tools %v\n", paastaversion.PaastaVersion)
		return 0, nil
	} else if err := cmd.Run(); err != nil {
		spanExec.Tag("error", err.Error())
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), nil
		}
		return 1, fmt.Errorf("error running %s: %s", subcommandPath, err)
	}
	return 0, nil
}

// os.Exit doesn't work well with defered calls
func main() {
	if len(os.Args) > 1 && os.Args[1] == "-version" {
		fmt.Printf("paasta-tools-go version: %v\n", paastaversion.Version)
		fmt.Printf("paasta-tools version: %v\n", paastaversion.PaastaVersion)
		fmt.Printf("zipkin initializers: %v\n", strings.Join(paastazipkin.Initializers(), ", "))
		fmt.Printf("go runtime: %v\n", runtime.Version())
		os.Exit(0)
	}

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	debug, _ := os.LookupEnv("PAASTA_DEBUG")
	v := klogFlags.Lookup("v")
	if v != nil {
		if debug != "" {
			v.Value.Set("10")
		} else {
			v.Value.Set("0")
		}
	}

	exit, err := paasta()
	if err != nil {
		klog.V(10).Infof("%v\n", err.Error())
	}
	os.Exit(exit)
}
