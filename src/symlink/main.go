package main

import (
	helpers "dotfiles/src"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func generateSymlink(source string, target string) {
	fmt.Printf("Symlinking: %s -> %s\n", source, target)

	if _, err := os.Stat(source); os.IsNotExist(err) {
		println("UNEXPECTED: Source not found:", source)
		return
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		println("EXPECTED: Target found, deleting:", target)

		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			println("UNEXPECTED: Error deleting target:", target)
			return
		}
	}

	targetDir := filepath.Dir(target)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		println("EXPECTED: Target directory not found, creating:", targetDir)

		mkdirErr := os.MkdirAll(targetDir, 0755)
		if mkdirErr != nil {
			println("UNEXPECTED: Error creating target directory:", targetDir)
			return
		}
	}

	err := os.Symlink(source, target)
	if err != nil {
		println("UNEXPECTED: Error creating symlink", err)
		return
	}

	println(source, "->", target)
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
		generateSymlink(sourcePath, targetPath)
	}

	helpers.PressAnyKeyOrWaitToExit()
}
