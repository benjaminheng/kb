package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/benjaminheng/kb/config"
	"github.com/fatih/color"
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
	// validate
	if err := runShellCommand("command -v ctags", nil, nil); err != nil {
		return errors.New("ctags not installed")
	}
	if config.Config.General.Editor != "vim" && config.Config.General.Editor != "nvim" {
		return errors.New("only vim is supported for this command. current editor: " + config.Config.General.Editor)
	}

	// generate tags
	tagsBuf := &bytes.Buffer{}
	command := "ctags --recurse=yes --with-list-header=no --machinable=yes -f - ."
	if err := runShellCommandWithWorkingDir(command, nil, tagsBuf, config.Config.General.KnowledgeBaseDir); err != nil {
		return err
	}

	var formattedTags string
	for _, line := range strings.Split(tagsBuf.String(), "\n") {
		components := strings.Split(line, "\t")
		if len(components) < 2 {
			continue
		}
		if config.Config.General.Color {
			components[0] = color.GreenString(components[0])
			components[1] = color.RedString(components[1])
		}
		// TODO: alignment
		formattedTag := strings.Join(components, "\t")
		formattedTags += formattedTag + "\n"
	}
	formattedTags = strings.TrimSpace(formattedTags)

	// get selection
	selectionBuf := &bytes.Buffer{}
	runSelectCommand(strings.NewReader(formattedTags), selectionBuf)
	if strings.TrimSpace(selectionBuf.String()) == "" {
		return nil
	}
	components := strings.Split(selectionBuf.String(), "\t")
	if len(components) == 0 {
		return nil
	}
	tag := components[0]

	// open in editor; editor has been validated to be either vim/nvim.
	command = fmt.Sprintf("%s -t \"%s\"", config.Config.General.Editor, tag)
	if err := runShellCommandWithWorkingDir(command, os.Stdin, os.Stdout, config.Config.General.KnowledgeBaseDir); err != nil {
		return err
	}
	return nil
}
