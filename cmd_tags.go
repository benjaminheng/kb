package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/benjaminheng/kb/config"
	"github.com/spf13/cobra"
)

func NewTagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags",
		Short: "Browse knowledge base by tag",
		Long:  ``,
		RunE:  tags,
	}
}

func tags(cmd *cobra.Command, args []string) error {
	// TODO: check if ctags exists
	// TODO: check if editor is vim

	// generate tags
	tagsBuf := &bytes.Buffer{}
	command := "ctags --recurse=yes --with-list-header=no --machinable=yes -f - ."
	if err := runShellCommandWithWorkingDir(command, nil, tagsBuf, config.Config.General.KnowledgeBaseDir); err != nil {
		return err
	}

	// get selection
	selectionBuf := &bytes.Buffer{}
	runShellCommand("fzf", tagsBuf, selectionBuf)
	if strings.TrimSpace(selectionBuf.String()) == "" {
		return nil
	}
	components := strings.Split(selectionBuf.String(), "\t")
	if len(components) == 0 {
		return nil
	}
	tag := components[0]

	// only vim is supported here
	command = fmt.Sprintf("%s -t \"%s\"", config.Config.General.Editor, tag)
	if err := runShellCommandWithWorkingDir(command, os.Stdin, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
		return err
	}
	return nil
}
