package main

import (
	"os/exec"
)

// Run starts the runs the command "cmd" in the "shell" after setting the
// environment variables to "env".
func Run(env []string, shell, cmd string) {
	command := exec.Command(shell)
	command.Env = env

	// We need to write to Stdin
	stdin, err := command.StdinPipe()
	if err != nil {
		return
	}

	command.Start()
	stdin.Write([]byte(cmd))
	stdin.Close()
}
