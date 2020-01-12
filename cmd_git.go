package main

import (
	"bufio"
	"fmt"
	"os"

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

func NewGitPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Pushes changes to remote",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: add, commit, push? Or is this taking too much
			// control away from the user?
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Commit message: ")
			text, _ := reader.ReadString('\n')
			fmt.Printf("text = %+v\n", text)
			// command := fmt.Sprintf("git commit -m \"%s\"", text)
			// if err := runShellCommandWithWorkingDir(command, nil, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
			// 	return err
			// }
			return nil
		},
	}
}
