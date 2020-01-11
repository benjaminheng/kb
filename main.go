package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/benjaminheng/kb/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	greenColor = color.New(color.FgGreen)
	redColor   = color.New(color.FgRed)
)

type rootCmdConfig struct {
}

type fileInfo struct {
	path  string
	title string
}

func browse(args []string) error {
	// browse files, display titles
	fileInfos := make([]fileInfo, 0)
	err := filepath.Walk(config.Config.General.KnowledgeBaseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fInfo := fileInfo{path: path}
			if config.Config.General.HasYAMLFrontMatter {
				f, err := os.Open(path)
				if err != nil {
					return err
				}
				result, err := ParseYAMLFrontMatter(f)
				if err != nil {
					return err
				}
				if v, ok := result["title"]; ok {
					if title, ok := v.(string); ok {
						fInfo.title = title
					}
				}
			}
			fileInfos = append(fileInfos, fInfo)
		}
		return nil
	})
	if err != nil {
		return err
	}
	var content string
	lineLookup := make(map[string]fileInfo)
	for i, v := range fileInfos {
		line := greenColor.Sprint(v.path)
		if v.title != "" {
			line += " :: " + redColor.Sprint(v.title)
		}
		lineLookup[line] = v
		content += line
		if i < len(fileInfos)-1 {
			content += "\n"
		}
	}
	r := strings.NewReader(content)
	b := &bytes.Buffer{}
	runShellCommand("fzf", r, b)
	selection := strings.Trim(b.String(), "\n")
	if selection != "" {
		if fInfo, ok := lineLookup[selection]; ok {
			editFile(config.Config.General.Editor, fInfo.path)
		}
	}
	return nil
}

func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "kb",
		Short: "Client for managing your knowledge base",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return browse(args)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := config.Config.Load()
			if err != nil {
				return err
			}
			if !config.Config.General.Color {
				redColor.DisableColor()
				greenColor.DisableColor()
			}
			return nil
		},
	}

	c.AddCommand(NewBrowseCmd())
	c.AddCommand(NewSearchCmd())

	// c.PersistentFlags().BoolVar(&rootConfig.Staging, "stage", false, "Use staging config")
	return c
}

func Execute() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
