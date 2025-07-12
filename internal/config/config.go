package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	GitHubTokenPath string            `yaml:"github_token_path"`
	Token           string            `yaml:"token"`
	SavedFilterSets []FilterSet       `yaml:"saved_filter_sets"`
	BaseBranches    []string          `yaml:"base_branches"`
	Theme           string            `yaml:"theme"`
	KeyBindings     map[string]string `yaml:"key_bindings"`
}

type FilterSet struct {
	Name   string   `yaml:"name"`
	Filter []string `yaml:"filter"`
}

func DefaultConfig() *Config {
	return &Config{
		GitHubTokenPath: "~/.github-token",
		BaseBranches:    []string{"main", "master", "develop"},
		Theme:           "default",
		KeyBindings:     make(map[string]string),
		SavedFilterSets: []FilterSet{
			{
				Name:   "Stale branches",
				Filter: []string{"STALE_LOCAL"},
			},
			{
				Name:   "Has PR",
				Filter: []string{"OPEN_PR", "DRAFT_PR", "CLOSED_PR"},
			},
		},
	}
}

func GetConfigPath() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	appConfigDir := filepath.Join(configDir, "branch-wrangler")
	if err := os.MkdirAll(appConfigDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, "config.yml"), nil
}

func Load() (*Config, error) {
	return DefaultConfig(), nil
}

func (c *Config) Save() error {
	return nil
}