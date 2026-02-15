package main

import (
	constants "dotfiles/src/constants"
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourceDir := filepath.Join(cwd, constants.SCRIPTS_SOURCE_DIR)
	outputDir := filepath.Join(cwd, constants.BUILD_SCRIPTS_DIR)

	if !utils.IsFileExists(outputDir) {
		os.MkdirAll(outputDir, 0755)
	}

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		buildScript(sourceDir, outputDir, entryName, entryName)

		aliasName := constants.BIN_SCRIPTS[entryName].Exe
		if aliasName != "" {
			buildScript(sourceDir, outputDir, entryName, aliasName)
		}
	}
}

func buildScript(sourceDir string, outputDir string, entryName string, exe string) {
	sourcePath := filepath.Join(sourceDir, entryName, "main.go")
	if !utils.IsFileExists(sourcePath) {
		panic(fmt.Sprintf("Source file not found: %s", sourcePath))
	}

	fmt.Println(aurora.Faint("> Building with Go: ").String() + entryName + aurora.Faint(" -> ").String() + exe)
	helpers.ExecNativeCommand([]string{"go", "build", "-o", filepath.Join(outputDir, exe+".exe"), sourcePath})
}
