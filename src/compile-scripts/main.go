package main

import (
	constants "dotfiles/src/constants"
	"dotfiles/src/helpers"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	scriptsDir := filepath.Join(cwd, constants.SOURCE_DIR_SCRIPTS)
	buildDir := filepath.Join(cwd, constants.BUILD_DIR_SCRIPTS)

	if !helpers.IsFileExists(buildDir) {
		os.MkdirAll(buildDir, 0755)
	}

	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		compileScript(entry.Name(), scriptsDir, buildDir)
	}
}

func compileScript(script string, sourceDir string, outputDir string) {
	outputPath := filepath.Join(outputDir, script+".exe")

	sourcePath := filepath.Join(sourceDir, script, "main.go")
	if !helpers.IsFileExists(sourcePath) {
		println("Source file not found", sourcePath)
		return
	}

	println("Building with Go", sourcePath, "to", outputPath)

	cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
