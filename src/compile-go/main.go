package main

import (
	constants "dotfiles/src/constants"
	"dotfiles/src/helpers"
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

	if !helpers.IsFileExists(outputDir) {
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
		outputName := constants.SCRIPTS_MAP[entryName].Exe
		if outputName == "" {
			outputName = entryName
		}

		sourcePath := filepath.Join(sourceDir, entryName, "main.go")
		outputPath := filepath.Join(outputDir, outputName+".exe")

		if !helpers.IsFileExists(sourcePath) {
			fmt.Println(aurora.Red("Source file not found: " + sourcePath))
			continue
		}

		fmt.Println(aurora.Faint("> Building with Go: ").String() + entryName)
		helpers.ExecNativeCommand([]string{"go", "build", "-o", outputPath, sourcePath})
	}
}
