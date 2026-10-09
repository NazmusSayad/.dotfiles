package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	constants "dotfiles/src/constants"
	"dotfiles/src/helpers"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	"github.com/otiai10/copy"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	sourceDir := filepath.Join(cwd, constants.SCRIPTS_SOURCE_DIR)
	var outputDir string
	switch runtime.GOOS {
	case "windows":
		outputDir = filepath.Join(cwd, ".local", "bin.win")
	case "darwin":
		outputDir = filepath.Join(cwd, ".local", "bin.mac")
	case "linux":
		outputDir = filepath.Join(cwd, ".local", "bin.wsl")
	default:
		panic("unsupported platform: " + runtime.GOOS)
	}

	if err := os.RemoveAll(outputDir); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		panic(err)
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
		scriptName := entryName
		if strings.HasSuffix(entryName, ".win") {
			if runtime.GOOS != "windows" {
				continue
			}
			scriptName = strings.TrimSuffix(entryName, ".win")
		} else if strings.HasSuffix(entryName, ".mac") {
			if runtime.GOOS != "darwin" {
				continue
			}
			scriptName = strings.TrimSuffix(entryName, ".mac")
		}

		if scriptName == "" {
			continue
		}

		aliasName := constants.BIN_SCRIPTS[scriptName].Exe
		if aliasName != "" {
			buildScript(sourceDir, outputDir, entryName, aliasName)
		} else {
			buildScript(sourceDir, outputDir, entryName, scriptName)
		}
	}

	fmt.Println(aurora.Faint("> Copying etc/bin -> ").String() + outputDir)
	if err := copy.Copy(filepath.Join(cwd, "etc", "bin"), outputDir); err != nil {
		panic(err)
	}
}

func buildScript(sourceDir string, outputDir string, entryName string, exe string) {
	sourcePath := filepath.Join(sourceDir, entryName, "main.go")
	if !utils.IsFileExists(sourcePath) {
		panic(fmt.Sprintf("Source file not found: %s", sourcePath))
	}

	binName := exe
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	fmt.Println(aurora.Faint("> Building with Go: ").String() + entryName + aurora.Faint(" -> ").String() + binName)
	helpers.ExecNativeCommand([]string{"go", "build", "-o", filepath.Join(outputDir, binName), sourcePath})
}
