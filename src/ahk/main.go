package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	helpers "dotfiles/src"
)

const ahkScriptPrefix = "AHK-"

func main() {
	baseDir := helpers.ResolvePath("./src/ahk")
	outputDir := helpers.ResolvePath("./bin/ahk")
	ahkScriptsDir := filepath.Join(baseDir, "scripts")
	ahk2ExeBin := filepath.Join(baseDir, "bin", "Ahk2Exe.exe")
	ahkCompilerBin := filepath.Join(baseDir, "bin", "AutoHotkey64.exe")

	entries, err := os.ReadDir(ahkScriptsDir)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
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
		outPath := filepath.Join(outputDir, ahkScriptPrefix+fileName+".exe")

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
