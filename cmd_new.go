package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
		f, err := os.Create(filePath)
		if err != nil {
			return err
		}
		io.WriteString(f, fmt.Sprintf("---\ntitle: \"%s\"\n---", title))
		f.Close()
		editFileWithWorkingDir(config.Config.General.Editor, filePath, config.Config.General.KnowledgeBaseDir)
	} else {
		editFileWithWorkingDir(config.Config.General.Editor, filePath, config.Config.General.KnowledgeBaseDir)
	}
	return nil
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
