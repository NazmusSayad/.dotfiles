package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	helpers "dotfiles/src/helpers"
)

type StringOrArray []string

func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}

	var array []string
	if err := json.Unmarshal(data, &array); err == nil {
		*s = array
		return nil
	}

	return fmt.Errorf("value must be a string or array of strings: %s", string(data))
}

type SymlinkConfig struct {
	Source  string
	Targets []string
}

type rawSymlinkConfig struct {
	Source    string        `json:"Source"`
	Target    StringOrArray `json:"Target"`
	TargetWin StringOrArray `json:"Target.Win"`
	TargetMac StringOrArray `json:"Target.Mac"`
}

func readSymlinkConfig() []SymlinkConfig {
	rawConfigs := helpers.ReadConfig[[]rawSymlinkConfig]("@/config/symlink.jsonc")

	var configs []SymlinkConfig
	for _, raw := range rawConfigs {
		targets := resolveTargets(raw)
		if len(targets) == 0 {
			continue
		}

		configs = append(configs, SymlinkConfig{
			Source:  raw.Source,
			Targets: targets,
		})
	}

	return configs
}

func resolveTargets(raw rawSymlinkConfig) []string {
	if runtime.GOOS == "windows" && len(raw.TargetWin) > 0 {
		return raw.TargetWin
	}

	if runtime.GOOS == "darwin" && len(raw.TargetMac) > 0 {
		return raw.TargetMac
	}

	if len(raw.Target) > 0 {
		return raw.Target
	}

	return nil
}

func main() {
	helpers.EnsureAdminExecution()
	symlinkConfigs := readSymlinkConfig()

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		if len(config.Targets) == 0 {
			fmt.Println("No targets found for", config.Source)
			continue
		}

		sourcePath := helpers.ResolvePath(config.Source)

		for _, target := range config.Targets {
			targetPath := helpers.ResolvePath(target)
			helpers.GenerateSymlink(sourcePath, targetPath)
		}
	}
}
