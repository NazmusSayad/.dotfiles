package main

import (
	constants "dotfiles/src/constants"
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourceDir := filepath.Join(cwd, constants.SCRIPTS_SOURCE_DIR)
	outputDir := filepath.Join(cwd, constants.SCRIPTS_BUILD_BIN_DIR)

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

		compileScript(entry.Name(), sourceDir, outputDir)
	}
}

func compileScript(script string, sourceDir string, outputDir string) {
	outputPath := filepath.Join(outputDir, script+".exe")

	sourcePath := filepath.Join(sourceDir, script, "main.go")
	if !helpers.IsFileExists(sourcePath) {
		fmt.Println("Source file not found", sourcePath)
		return
	}

	fmt.Println("Building with Go", sourcePath, "to", outputPath)

	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
