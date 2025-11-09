package winget

import (
	helpers "dotfiles/src"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type WingetPackage struct {
	ID      string
	Name    string
	Version string

	SilentInstall bool
	SilentUpgrade bool

	IgnoreInstall bool
	IgnoreUpgrade bool
}

func GetWingetPackages(path string) []WingetPackage {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		os.Exit(1)
	}

	fmt.Printf("CWD: %s\n", cwd)

	fullPath := filepath.Join(cwd, path)
	jsonBytes, err := helpers.ReadJSONWithComments(fullPath)
	if err != nil {
		return []WingetPackage{}
	}

	var pkgs []WingetPackage
	if err := json.Unmarshal(jsonBytes, &pkgs); err != nil {
		return []WingetPackage{}
	}

	return pkgs
}

func BuildWingetInstallCommands(p WingetPackage) []string {
	parts := []string{"winget", "install", p.ID}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	if p.SilentInstall {
		parts = append(parts, "--silent")
	} else {
		parts = append(parts, "--interactive")
	}

	return append(parts, "--verbose", "--accept-package-agreements", "--accept-source-agreements", "--no-upgrade")
}

func BuildWingetUpgradeCommands(p WingetPackage) []string {
	parts := []string{"winget", "upgrade", p.ID}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	if p.SilentUpgrade {
		parts = append(parts, "--silent")
	} else {
		parts = append(parts, "--interactive")
	}

	return append(parts, "--verbose", "--accept-package-agreements", "--accept-source-agreements")
}
