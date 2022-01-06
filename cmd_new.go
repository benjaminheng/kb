package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/benjaminheng/kb/config"
	"github.com/spf13/cobra"
)

func NewNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new file",
		Long:  ``,
		RunE:  createNewFile,
	}
	return cmd
}

func createNewFile(cmd *cobra.Command, args []string) error {
	if config.Config.General.Editor != "vim" && config.Config.General.Editor != "nvim" {
		// Currently only vim is supported, because we'll be opening
		// vim with an unsaved buffer. For other editors, we might have
		// to write the file first, then open it in the editor.
		return errors.New("`kb new` is only supported if editor=vim or editor=nvim in config")
	}

	scanner := bufio.NewScanner(os.Stdin)

	var title, filename string

	fmt.Print("Title: ")
	if scanner.Scan() {
		title = scanner.Text()
		if title == "" {
			return errors.New("title required")
		}
	}

	filename = titleToFilename(title)
	fmt.Printf("Filename (%s): ", filename)
	if scanner.Scan() {
		s := scanner.Text()
		if s != "" {
			filename = s
		}
	}

	filePath := path.Join(config.Config.General.KnowledgeBaseDir, filename)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return editUnsavedFileInVim(config.Config.General.Editor, filePath, strings.NewReader(fmt.Sprintf("---\ntitle: \"%s\"\n---", title)), os.Stdout, config.Config.General.KnowledgeBaseDir)

	}
	return editFileWithWorkingDir(config.Config.General.Editor, filePath, config.Config.General.KnowledgeBaseDir)
}

func titleToFilename(title string) string {
	nonWords := regexp.MustCompile(`\W`)
	consecutiveDashes := regexp.MustCompile(`-{2,}`)
	filename := nonWords.ReplaceAllString(title, "-")
	filename = consecutiveDashes.ReplaceAllString(filename, "-")
	filename = strings.ToLower(filename)
	filename += ".md"
	return filename
}

// validate that file exists and is not a directory
func validateFileExists(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", filePath)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory", filePath)
	}
	return nil
}
