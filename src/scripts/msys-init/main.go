package main

import (
	helpers "dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	helpers.EnsureAdminExecution()

	const MSYS_PATH = "C:\\msys64"
	NSSWITCH_CONFIG_PATH := filepath.Join(MSYS_PATH, "etc", "nsswitch.conf")

	MSYS_INI_LIST := []string{
		"msys2.ini",
		"ucrt64.ini",

		"mingw32.ini",
		"mingw64.ini",

		"clang64.ini",
		"clangarm64.ini",
	}

	MSYS_BIN_LIST := []string{
		"usr/bin",
		"ucrt64/bin",

		"mingw32/bin",
		"mingw64/bin",

		"clang64/bin",
		"clangarm64/bin",
	}

	pathTypeRegex := regexp.MustCompile(`(?m)^#MSYS2_PATH_TYPE=inherit`)
	msysWinSymlinksRegex := regexp.MustCompile(`(?m)^#MSYS=winsymlinks:nativestrict`)

	for _, ini := range MSYS_INI_LIST {
		iniPath := filepath.Join(MSYS_PATH, ini)
		if !utils.IsFileExists(iniPath) {
			fmt.Println("File not found:", iniPath)
			continue
		}

		content, err := os.ReadFile(iniPath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}

		updated := pathTypeRegex.ReplaceAll(content, []byte("MSYS2_PATH_TYPE=inherit"))
		updated = msysWinSymlinksRegex.ReplaceAll(updated, []byte("MSYS=winsymlinks:nativestrict"))

		_ = os.WriteFile(iniPath, updated, 0644)
		fmt.Println("Updated:", ini)
	}

	if content, err := os.ReadFile(NSSWITCH_CONFIG_PATH); err == nil {
		reNss := regexp.MustCompile(`(?m)^(db_home|db_shell|db_gecos):\s*.*$`)
		updated := reNss.ReplaceAllString(string(content), "$1: windows")
		_ = os.WriteFile(NSSWITCH_CONFIG_PATH, []byte(updated), 0644)
		fmt.Println("Updated: nsswitch.conf")
	}

	fmt.Println("Adding MSYS2_PATH_TYPE to env")
	helpers.WriteEnv(helpers.ScopeMachine, "MSYS2_PATH_TYPE", "inherit")

	var existingBins []string
	for _, rel := range MSYS_BIN_LIST {
		full := filepath.Join(MSYS_PATH, rel)

		if info, err := os.Stat(full); err == nil && info.IsDir() {
			existingBins = append(existingBins, full)
		}
	}

	fmt.Println("Adding to PATH:", len(existingBins), "bins")
	_, _ = helpers.AddToEnvPath(helpers.ScopeMachine, existingBins...)
}
