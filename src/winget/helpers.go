package winget

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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
	exePath, e := os.Executable()
	if e != nil {
		os.Exit(1)
	}

	fullPath := filepath.Join(filepath.Dir(exePath), path)
	fmt.Println("Loading winget packages from", fullPath)

	f, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("File opening failed...")
		return []WingetPackage{}
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("File reading failed...")
		return []WingetPackage{}
	}

	fmt.Println("File loaded successfully...")
	re := regexp.MustCompile(`(?m)//.*$`)
	clean := re.ReplaceAll(data, []byte{})

	reBlock := regexp.MustCompile(`(?s)/\*.*?\*/`)
	clean = reBlock.ReplaceAll(clean, []byte{})

	var pkgs []WingetPackage
	if err := json.Unmarshal(clean, &pkgs); err != nil {
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
