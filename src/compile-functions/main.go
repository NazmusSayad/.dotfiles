package main

import (
	"dotfiles/src/helpers"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	constants "dotfiles/src"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	functionsDir := filepath.Join(cwd, constants.SOURCE_DIR_FUNCTIONS)
	buildDir := filepath.Join(cwd, constants.BUILD_DIR_FUNCTIONS)

	if !helpers.IsFileExists(buildDir) {
		os.MkdirAll(buildDir, 0755)
	}

	entries, err := os.ReadDir(functionsDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".rs") {
			continue
		}

		fileName := strings.TrimSuffix(entry.Name(), ".rs")
		sourcePath := filepath.Join(functionsDir, entry.Name())
		outputPath := filepath.Join(buildDir, fileName+".exe")

		println("Compiling", fileName, "from", sourcePath, "to", outputPath)

		cmd := exec.Command("rustc", "-C", "strip=symbols", "-Clink-arg=/DEBUG:NONE", "-Clink-arg=/PDB:NONE", sourcePath, "-o", outputPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
