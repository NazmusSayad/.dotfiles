package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	constants "dotfiles/src/constants"
	helpers "dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

const ahkScriptPrefix = "AHK-"

func main() {
	if err := downloadFromReleases(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	buildOutputDir := filepath.Join(constants.BUILD_DIR, "/ahk")
	ahkScriptsDir := filepath.Join(constants.SOURCE_DIR, "/compile-ahk/scripts")

	ahk2ExeBin := filepath.Join(constants.BUILD_LIBRARIES_DIR, "Ahk2Exe.exe")
	ahk64CompilerBin := filepath.Join(constants.BUILD_LIBRARIES_DIR, "AutoHotkey64.exe")

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
		spawnArgs := []string{"/base", ahk64CompilerBin, "/in", inPath, "/out", outPath}
		if _, err := os.Stat(iconPath); err == nil {
			spawnArgs = append(spawnArgs, "/icon", iconPath)
		}

		err := helpers.ExecNativeCommand(append([]string{ahk2ExeBin}, spawnArgs...))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println(aurora.Faint("> Compiled: ").String() + entry.Name())
	}
}

func downloadFromReleases() error {
	err := helpers.WriteGithubReleaseZipFile(
		constants.BUILD_LIBRARIES_DIR,
		"https://github.com/AutoHotkey/AutoHotkey",
		"AutoHotkey_.*",
		"AutoHotkey64.exe",
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	err2 := helpers.WriteGithubReleaseZipFile(
		constants.BUILD_LIBRARIES_DIR,
		"https://github.com/AutoHotkey/Ahk2Exe",
		"Ahk2Exe.*",
		"Ahk2Exe.exe",
	)

	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
		return err2
	}

	err3 := helpers.WriteGithubReleaseFile(
		constants.BUILD_LIBRARIES_DIR,
		"https://github.com/Ciantic/VirtualDesktopAccessor",
		"VirtualDesktopAccessor.dll",
	)

	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
		return err3
	}

	return nil
}
