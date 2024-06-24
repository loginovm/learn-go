package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	return RunCmdWithWriters(cmd, env, os.Stdout, os.Stderr)
}

func RunCmdWithWriters(cmd []string, env Environment, stdout, stderr io.Writer) (returnCode int) {
	command := cmd[0]
	args := cmd[1:]
	c := exec.Command(command, args...)

	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v.Value)
		}
	}

	c.Stdout = stdout
	c.Stderr = stderr

	if err := c.Run(); err != nil {
		stderr.Write([]byte(err.Error()))
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}

	return 0
}
