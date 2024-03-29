package main

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/benjaminheng/kb/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type fileInfo struct {
	path     string
	filename string
	title    string
}

func browse(args []string) error {
	// browse files, display titles
	fileInfos := make([]fileInfo, 0)
	err := filepath.Walk(config.Config.General.KnowledgeBaseDir, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			fInfo := fileInfo{path: fpath, filename: path.Base(fpath)}
			if config.Config.IsIgnoredFilename(fInfo.filename) {
				return nil
			}
			if config.Config.General.HasYAMLFrontMatter {
				f, err := os.Open(fpath)
				if err != nil {
					return err
				}
				result, err := ParseYAMLFrontMatter(f)
				if err != nil && err != ErrFrontMatterNotFound {
					return err
				}
				if err == nil {
					if v, ok := result["title"]; ok {
						if title, ok := v.(string); ok {
							fInfo.title = title
						}
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
		var line string
		if v.title != "" {
			line = v.title + " :: "
		}
		line += v.path
		lineLookup[line] = v

		if config.Config.General.Color {
			if v.title != "" {
				line = color.GreenString(v.title) + " :: "
			}
			line += v.path
		}

		content += line
		if i < len(fileInfos)-1 {
			content += "\n"
		}
	}
	r := strings.NewReader(content)
	b := &bytes.Buffer{}
	runSelectCommand(r, b)
	selection := strings.Trim(b.String(), "\n")
	if selection != "" {
		if fInfo, ok := lineLookup[selection]; ok {
			editFileWithWorkingDir(config.Config.General.Editor, fInfo.path, config.Config.General.KnowledgeBaseDir)
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
			err := config.Config.Load(config.Config.General.ConfigFile)
			if err != nil {
				return err
			}
			return nil
		},
	}

	c.AddCommand(NewBrowseCmd())
	c.AddCommand(NewSearchCmd())
	c.AddCommand(NewGitCmd())
	c.AddCommand(NewNewCmd())

	c.PersistentFlags().StringVar(&config.Config.General.ConfigFile, "config", "", "Path to config file (default: $HOME/.config/kb/config.toml)")
	return c
}
