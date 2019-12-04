package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
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

func main() {
	var cmds map[string]bool = nil
	var cmdPath string
	var args []string

	if len(os.Args) > 1 {
		var err error
		cmds, err = listPaastaCommands()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Error generating list of sub-commands: %s\n",
				err,
			)
			os.Exit(1)
		}

		fullCmdName := fmt.Sprintf("paasta-%s", os.Args[1])
		if _, ok := cmds[fullCmdName]; ok {
			var err error
			cmdPath, err = exec.LookPath(fullCmdName)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Couldn't lookup %s in PATH: %s\n",
					fullCmdName,
					err,
				)
				os.Exit(1)
			}
		}
	}

	if cmdPath == "" {
		cmdPath = "/opt/venvs/paasta-tools/bin/paasta"
		args = os.Args
	} else {
		args = os.Args[1:]
	}

	if err := syscall.Exec(cmdPath, args, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "Error running %s: %s", cmdPath, err)
		os.Exit(1)
	}
}
