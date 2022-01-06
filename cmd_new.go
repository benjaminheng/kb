package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"

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
	return nil
}

func titleToFilename(title string) string {
	nonWords := regexp.MustCompile(`\W`)
	consecutiveDashes := regexp.MustCompile(`-{2,}`)
	filename := nonWords.ReplaceAllString(title, "-")
	filename = consecutiveDashes.ReplaceAllString(filename, "-")
	filename += ".md"
	return filename
}
