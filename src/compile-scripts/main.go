package main

import (
	"dotfiles/src/helpers"
	"os"
	"os/exec"
	"path/filepath"
)

const SCRIPTS_DIR = "./scripts"
const BUILD_DIR = "./build/bin"

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	scriptsDir := filepath.Join(cwd, SCRIPTS_DIR)
	buildDir := filepath.Join(cwd, BUILD_DIR)

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

		packageName := entry.Name()
		sourcePath := filepath.Join(scriptsDir, packageName, "main.go")
		outputPath := filepath.Join(buildDir, packageName+".exe")

		println("Compiling", packageName, "from", sourcePath, "to", outputPath)

		cmd := exec.Command("go", "build", "-o", outputPath, sourcePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
