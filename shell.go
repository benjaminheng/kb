package main

import (
	"io"
	"os"
	"os/exec"
	"runtime"
)

func editFile(command, file string) error {
	command += " " + file
	return runShellCommand(command, os.Stdin, os.Stdout)
}

func runShellCommand(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}
