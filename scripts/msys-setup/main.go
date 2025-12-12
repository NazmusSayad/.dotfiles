package main

import (
	helpers "dotfiles/src/helpers"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	helpers.EnsureAdminExecution()

	const MSYS_PATH = "C:\\msys64"
	NSSWITCH_CONFIG_PATH := filepath.Join(MSYS_PATH, "etc", "nsswitch.conf")

	MSYS_INIS := []string{
		"msys2.ini",
		"clang32.ini",
		"clang64.ini",
		"clangarm64.ini",
		"mingw32.ini",
		"mingw64.ini",
		"ucrt64.ini",
	}

	MSYS_BINS := []string{
		"usr/bin",
		"mingw64/bin",
		"mingw32/bin",
		"ucrt64/bin",
		"clang32/bin",
		"clang64/bin",
		"clangarm64/bin",
	}

	reIni := regexp.MustCompile(`(?m)^#MSYS2_PATH_TYPE=inherit`)
	for _, ini := range MSYS_INIS {
		iniPath := filepath.Join(MSYS_PATH, ini)
		if !helpers.IsFileExists(iniPath) {
			println("File not found: %s\n", iniPath)
			continue
		}

		content, err := os.ReadFile(iniPath)
		if err != nil {
			continue
		}

		updated := reIni.ReplaceAll(content, []byte("MSYS2_PATH_TYPE=inherit"))
		_ = os.WriteFile(iniPath, updated, 0644)
		println("Updated: %s\n", ini)
	}

	if content, err := os.ReadFile(NSSWITCH_CONFIG_PATH); err == nil {
		reNss := regexp.MustCompile(`(?m)^(db_home|db_shell|db_gecos):\s*.*$`)
		updated := reNss.ReplaceAllString(string(content), "$1: windows")
		_ = os.WriteFile(NSSWITCH_CONFIG_PATH, []byte(updated), 0644)
		println("Updated: nsswitch.conf")
	}

	_, _ = helpers.WriteEnv(helpers.ScopeMachine, "MSYS2_PATH_TYPE", "inherit")

	var existingBins []string
	for _, rel := range MSYS_BINS {
		full := filepath.Join(MSYS_PATH, rel)
		if info, err := os.Stat(full); err == nil && info.IsDir() {
			existingBins = append(existingBins, full)
		}
	}
	_, _ = helpers.AddToEnvPath(helpers.ScopeMachine, existingBins...)
}
