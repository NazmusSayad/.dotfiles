package main

import (
	helpers "dotfiles/src"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func generateSymlink(source string, target string) {
	fmt.Printf("Symlinking: %s -> %s\n", source, target)

	if _, err := os.Stat(source); os.IsNotExist(err) {
		fmt.Println("UNEXPECTED: Source not found:", source)
		return
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		fmt.Println("EXPECTED: Target found, deleting:", target)

		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println("UNEXPECTED: Error deleting target:", target)
			return
		}
	}

	targetDir := filepath.Dir(target)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Println("EXPECTED: Target directory not found, creating:", targetDir)

		mkdirErr := os.MkdirAll(targetDir, 0755)
		if mkdirErr != nil {
			fmt.Println("UNEXPECTED: Error creating target directory:", targetDir)
			return
		}
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println("UNEXPECTED: Error creating symlink", err)
		return
	}

	fmt.Println(source, "->", target)
}

func parseSymlinkConfig(path string) []SymlinkConfig {
	jsonBytes, err := helpers.ReadJsoncAsJson(path)
	if err != nil {
		fmt.Println("Error reading JSON file...")
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

	symlinkConfigPath := helpers.ResolvePath("./config/symlink.jsonc")
	symlinkConfigs := parseSymlinkConfig(symlinkConfigPath)

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		time.Sleep(2000)
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)
		targetPath := helpers.ResolvePath(config.Target)

		fmt.Println("")
		generateSymlink(sourcePath, targetPath)
	}

	helpers.PressAnyKeyOrWaitToExit()
}
