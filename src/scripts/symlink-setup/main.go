package main

import (
	helpers "dotfiles/src/helpers"
	"fmt"
	"os"
)

type SymlinkConfig struct {
	Source  string
	Target  string
	Targets []string
}

func main() {
	helpers.EnsureAdminExecution()
	symlinkConfigs := helpers.ReadConfig[[]SymlinkConfig]("@/config/symlink.jsonc")

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		targets := []string{}
		if config.Target != "" {
			targets = append(targets, config.Target)
		}
		if len(config.Targets) > 0 {
			targets = append(targets, config.Targets...)
		}

		if len(targets) == 0 {
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
