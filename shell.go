package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/benjaminheng/kb/config"
)

func editFile(command, file string) error {
	command += " " + file
	return runShellCommand(command, os.Stdin, os.Stdout)
}

func editFileWithWorkingDir(command, file, workingDir string) error {
	return editFile(withWorkingDir(workingDir, command), file)
}

func runShellCommandWithWorkingDir(command string, r io.Reader, w io.Writer, workingDir string) error {
	return runShellCommand(withWorkingDir(workingDir, command), r, w)
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

func withWorkingDir(workingDir, command string) string {
	return fmt.Sprintf("cd %s && %s", workingDir, command)
}

func runSelectCommand(r io.Reader, w io.Writer) error {
	command := config.Config.General.SelectCmd
	if command == "fzf" && config.Config.General.Color {
		command += " --ansi"
	}
	return runShellCommand(command, r, w)
}

func getUserInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	return text
}
