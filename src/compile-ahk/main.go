package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	constants "dotfiles/src/constants"
	helpers "dotfiles/src/helpers"
)

const ahkScriptPrefix = "AHK-"

func main() {
	srcDir := helpers.ResolvePath(constants.SOURCE_DIR + "/compile-ahk")
	buildOutputDir := helpers.ResolvePath(constants.BUILD_DIR + "/ahk")

	ahkScriptsDir := filepath.Join(srcDir, "scripts")
	ahk2ExeBin := filepath.Join(srcDir, "bin", "Ahk2Exe.exe")
	ahkCompilerBin := filepath.Join(srcDir, "bin", "AutoHotkey64.exe")

	entries, err := os.ReadDir(ahkScriptsDir)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(buildOutputDir, 0755); err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".ahk") {
			continue
		}

		fileName := strings.TrimSuffix(entry.Name(), ".ahk")
		inPath := filepath.Join(ahkScriptsDir, entry.Name())
		outPath := filepath.Join(buildOutputDir, ahkScriptPrefix+fileName+".exe")

		iconPath := filepath.Join(ahkScriptsDir, fileName+".ico")
		spawnArgs := []string{"/base", ahkCompilerBin, "/in", inPath, "/out", outPath}
		if _, err := os.Stat(iconPath); err == nil {
			spawnArgs = append(spawnArgs, "/icon", iconPath)
		}

		cmd := exec.Command(ahk2ExeBin, spawnArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		_ = cmd.Run()
		fmt.Printf("Compiled: %s\n", entry.Name())
	}

	println("Compilation complete")
}
