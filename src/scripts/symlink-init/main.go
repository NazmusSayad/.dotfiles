package main

import (
	"fmt"
	"os"
	"runtime"

	helpers "dotfiles/src/helpers"
)

type SymlinkConfig struct {
	Source  string
	Targets []string
}

type rawSymlinkConfig struct {
	Source     string   `json:"Source"`
	Target     string   `json:"Target"`
	Targets    []string `json:"Targets"`
	TargetWin  string   `json:"Target.Win"`
	TargetMac  string   `json:"Target.Mac"`
	TargetsWin []string `json:"Targets.Win"`
	TargetsMac []string `json:"Targets.Mac"`
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
	if runtime.GOOS == "windows" {
		if len(raw.TargetsWin) > 0 {
			return raw.TargetsWin
		}
		if raw.TargetWin != "" {
			return []string{raw.TargetWin}
		}
	}

	if runtime.GOOS == "darwin" {
		if len(raw.TargetsMac) > 0 {
			return raw.TargetsMac
		}
		if raw.TargetMac != "" {
			return []string{raw.TargetMac}
		}
	}

	if len(raw.Targets) > 0 {
		return raw.Targets
	}

	if raw.Target != "" {
		return []string{raw.Target}
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
		targets := []string{}

		if len(config.Targets) > 0 {
			targets = append(targets, config.Targets...)
		}

		if len(targets) == 0 {
			fmt.Println("No targets found for", config.Source)
			continue
		}

		sourcePath := helpers.ResolvePath(config.Source)

		for _, target := range targets {
			targetPath := helpers.ResolvePath(target)
			helpers.GenerateSymlink(sourcePath, targetPath)
		}
	}
}
