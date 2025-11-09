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

func parseSymlinkConfig(cwd string, path string) []SymlinkConfig {
	fullPath := filepath.Join(cwd, path)
	jsonBytes, err := helpers.ReadJSONWithComments(fullPath)
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

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		os.Exit(1)
	}

	fmt.Printf("CWD: %s\n", cwd)
	symlinkConfigs := parseSymlinkConfig(cwd, "./config/symlink.jsonc")

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(cwd, config.Source)
		targetPath := helpers.ResolvePath(cwd, config.Target)

		fmt.Println("")
		generateSymlink(sourcePath, targetPath)
	}

	helpers.WaitForInputAndExit()
}
