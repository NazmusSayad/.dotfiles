package agentstore

import (
	"path/filepath"
	"slices"

	"dotfiles/src/helpers"
)

const defaultConfigPath = "@/config/agents-store.yml"

type Config struct {
	Database  string   `yaml:"database"`
	ClaudeDir string   `yaml:"claudeDir"`
	Items     []string `yaml:"items"`
	Skip      []string `yaml:"skip"`
	NoPull    []string `yaml:"noPull"`
}

func LoadConfig(args []string) Config {
	if len(args) > 1 {
		fatal("Usage: <command> [config path]")
	}
	path := defaultConfigPath
	if len(args) == 1 {
		path = args[0]
	}

	config := helpers.ReadConfig[Config](path)
	switch {
	case config.Database == "":
		fatal("database is missing in " + path)
	case config.ClaudeDir == "":
		fatal("claudeDir is missing in " + path)
	case len(config.Items) == 0:
		fatal("items is missing in " + path)
	}
	for _, pattern := range config.Skip {
		if _, err := filepath.Match(pattern, ""); err != nil {
			fatal("invalid skip pattern " + pattern + " in " + path)
		}
	}

	config.ClaudeDir = helpers.ResolvePath(config.ClaudeDir)
	return config
}

func (c Config) skips(name string) bool {
	for _, pattern := range c.Skip {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return false
}

func (c Config) pulls(path string) bool {
	return !slices.Contains(c.NoPull, path)
}
