package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/benjaminheng/kb/config"
	"github.com/spf13/cobra"
)

func NewGitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Performs git operations on the knowledge base",
		Long:  ``,
	}
	cmd.AddCommand(NewGitPullCmd())
	cmd.AddCommand(NewGitPushCmd())
	cmd.AddCommand(NewGitStatusCmd())
	return cmd
}

func NewGitPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pulls changes from remote",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runShellCommandWithWorkingDir("git pull", nil, os.Stdout, config.Config.General.KnowledgeBaseDir)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func NewGitStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Shows working tree status",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runShellCommandWithWorkingDir("git status", nil, os.Stdout, config.Config.General.KnowledgeBaseDir)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

func NewGitPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Pushes changes to remote",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// check if there are files to commit, if not then do
			// `git status` and exit quietly
			if err := runShellCommandWithWorkingDir("git diff-index --quiet HEAD --", nil, os.Stdout, config.Config.General.KnowledgeBaseDir); err == nil {
				runShellCommandWithWorkingDir("git status", nil, os.Stdout, config.Config.General.KnowledgeBaseDir)
				return nil
			}

			// git add .
			fmt.Println("+ git add .")
			if err := runShellCommandWithWorkingDir("git add .", nil, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
				return err
			}

			// git status
			fmt.Println("+ git status")
			if err := runShellCommandWithWorkingDir("git status", nil, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
				return err
			}

			// ask user for confirmation
			text := getUserInput("Confirm (Y/n): ")
			if text != "" && strings.ToLower(text) != "y" {
				fmt.Println("Executing: git reset .")
				return runShellCommandWithWorkingDir("git reset .", nil, os.Stdout, config.Config.General.KnowledgeBaseDir)
			}

			// get list of changed files
			b := &bytes.Buffer{}
			if err := runShellCommandWithWorkingDir("git diff --name-only --cached", nil, b, config.Config.General.KnowledgeBaseDir); err != nil {
				return err
			}
			files := strings.Split(b.String(), "\n")
			if len(files) == 0 {
				return errors.New("no files staged for commit") // should not happen
			}
			defaultCommitMessage := fmt.Sprintf("Update %s", files[0])

			// ask user for commit message
			text = getUserInput(fmt.Sprintf("Commit message (%s): ", defaultCommitMessage))
			if text == "" {
				text = defaultCommitMessage
			}

			// git commit -m "X"
			command := fmt.Sprintf("git commit -m \"%s\"", text)
			fmt.Println("+ " + command)
			if err := runShellCommandWithWorkingDir(command, nil, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
				return err
			}

			// git push
			fmt.Println("+ git push")
			if err := runShellCommandWithWorkingDir("git push", nil, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
				return err
			}
			return nil
		},
	}
}
