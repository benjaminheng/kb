package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config is a global variable used to store general configuration.
var Config RootConfig

// RootConfig represents the configuration structure of the .toml configuration
// file.
type RootConfig struct {
	General GeneralConfig `toml:"general"`
}

// GeneralConfig represents general configuration values.
type GeneralConfig struct {
	KnowledgeBaseDir   string `toml:"knowledge_base_dir"`
	Editor             string `toml:"editor"`
	SelectCmd          string `toml:"select_cmd"`
	HasYAMLFrontMatter bool   `toml:"has_yaml_front_matter"`
	Color              bool   `toml:"color"`
}

// Flag is a global variable used to store flags.
var Flag FlagConfig

// FlagConfig represents flags that the kb command accepts.
type FlagConfig struct {
}

func (cfg *RootConfig) Load() error {
	file, err := getConfigFile()
	if err != nil {
		return err
	}

	// Create default config if it does not already exist
	_, err = os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.initDefaultConfig()
			return nil
		}
		return err
	}

	// If config exists, try to load from it
	_, err = toml.DecodeFile(file, cfg)
	if err != nil {
		return err
	}
	return nil
}

func (cfg *RootConfig) initDefaultConfig() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}
	configFile, err := os.Create(filepath.Join(configDir, "config.toml"))
	if err != nil {
		return err
	}
	kbDir := filepath.Join(configDir, "kb")
	if err := os.MkdirAll(kbDir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}

	cfg.General.KnowledgeBaseDir = kbDir
	cfg.General.SelectCmd = "fzf"
	cfg.General.Color = false

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		cfg.General.Editor = "vim"
	}

	return toml.NewEncoder(configFile).Encode(cfg)
}

func getConfigDir() (string, error) {
	dir := filepath.Join(os.Getenv("HOME"), ".config", "kb")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func getConfigFile() (configFile string, err error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	configFile = filepath.Join(dir, "config.toml")
	return configFile, nil
}
