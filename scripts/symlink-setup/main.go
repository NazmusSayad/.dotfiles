package main

import (
	helpers "dotfiles/src/helpers"
	"encoding/json"
	"os"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func parseSymlinkConfig() []SymlinkConfig {
	jsonBytes, err := helpers.ReadDotfilesConfigJSONC("./config/symlink.jsonc")

	if err != nil {
		println("Error reading JSON file...")
		return []SymlinkConfig{}
	}

	var symlinkConfigs []SymlinkConfig
	if err := json.Unmarshal(jsonBytes, &symlinkConfigs); err != nil {
		return []SymlinkConfig{}
	}

	return symlinkConfigs
}

func main() {
	helpers.EnsureAdminExecution()

	symlinkConfigs := parseSymlinkConfig()
	if len(symlinkConfigs) == 0 {
		println("No symlink configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)
		targetPath := helpers.ResolvePath(config.Target)

		println("")
		helpers.GenerateSymlink(sourcePath, targetPath)
	}
}
