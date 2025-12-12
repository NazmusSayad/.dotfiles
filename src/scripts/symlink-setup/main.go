package main

import (
	helpers "dotfiles/src/helpers"
	"encoding/json"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func parseSymlinkConfig() []SymlinkConfig {
	jsonBytes, err := helpers.ReadDotfilesConfigJSONC("./config/symlink.jsonc")

	if err != nil {
		fmt.Println(aurora.Red("Error reading JSON file..."))
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
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)
		targetPath := helpers.ResolvePath(config.Target)

		fmt.Println()
		helpers.GenerateSymlink(sourcePath, targetPath)
	}
}
