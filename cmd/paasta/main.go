package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func listPaastaCommands() (map[string]bool, error) {
	// alternatively:
	// find $(echo $PATH | tr ':' ' ') -maxdepth 1 -xtype f -perm /o+x -name paasta-*
	cmd := exec.Command("/bin/bash", "-p", "-c", "compgen -A command paasta-")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
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
	var subcmd string
	var args []string

	if len(os.Args) > 1 {
		var err error
		cmds, err = listPaastaCommands()
		if err != nil {
			fmt.Printf("Error generating list of sub-commands: %s\n", err)
			os.Exit(1)
		}

		cmd := fmt.Sprintf("paasta-%s", os.Args[1])
		if ok, _ := cmds[cmd]; ok {
			var err error
			subcmd, err = exec.LookPath(cmd)
			if err != nil {
				fmt.Printf("Couldn't lookup %s in PATH: %s\n", cmd, err)
				os.Exit(1)
			}
		}
	}

	if subcmd == "" {
		subcmd = "/usr/bin/paasta"
		args = os.Args
	} else {
		args = os.Args[1:]
	}

	if err := syscall.Exec(subcmd, args, os.Environ()); err != nil {
		fmt.Printf("Error running %s: %s", subcmd, err)
		os.Exit(1)
	}
}
