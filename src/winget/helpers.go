package winget

import (
	helpers "dotfiles/src"
	"encoding/json"
	"os/exec"
	"strings"
)

type WingetPackage struct {
	ID      string
	Name    string
	Version string

	InteractiveInstall bool
	InteractiveUpgrade bool

	SkipInstall bool
	SkipUpgrade bool
}

type WingetUpgradeablePackage struct {
	ID        string
	Version   string
	Available string
}

func GetWingetPackages() []WingetPackage {
	jsonBytes, err := helpers.ReadDotfilesConfigJSONC("./config/winget-apps.jsonc")
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

	if p.InteractiveInstall {
		parts = append(parts, "--interactive")
	} else {
		parts = append(parts, "--silent")
	}

	return append(parts, "--verbose", "--accept-package-agreements", "--accept-source-agreements", "--no-upgrade")
}

func BuildWingetUpgradeCommands(p WingetPackage) []string {
	parts := []string{"winget", "upgrade", p.ID}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	if p.InteractiveUpgrade {
		parts = append(parts, "--interactive")
	} else {
		parts = append(parts, "--silent")
	}

	return append(parts, "--verbose", "--accept-package-agreements", "--accept-source-agreements")
}

func GetUpgradeablePackages() []WingetUpgradeablePackage {
	var upgradeablePackages []WingetUpgradeablePackage

	cmd := exec.Command("winget", "upgrade")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return upgradeablePackages
	}

	trimmedOutput := strings.TrimSpace(string(output))
	cleanedOutput := strings.ReplaceAll(trimmedOutput, "\r\n", "\n")
	lines := strings.Split(cleanedOutput, "\n")

	headingLine := strings.TrimSpace(lines[0])
	if headingLine == "" {
		return upgradeablePackages
	}

	dataLines := lines[2 : len(lines)-1]
	if len(dataLines) == 0 {
		return upgradeablePackages
	}

	nameIndex := strings.Index(headingLine, "Name")
	idStartIndex := strings.Index(headingLine, "Id") - nameIndex
	versionStartIndex := strings.Index(headingLine, "Version") - nameIndex
	availableStartIndex := strings.Index(headingLine, "Available") - nameIndex
	sourceStartIndex := strings.Index(headingLine, "Source") - nameIndex

	for _, line := range dataLines {
		line = strings.TrimSpace(line)

		upgradeablePackages = append(upgradeablePackages, WingetUpgradeablePackage{
			ID:        strings.TrimSpace(line[idStartIndex:versionStartIndex]),
			Version:   strings.TrimSpace(line[versionStartIndex:availableStartIndex]),
			Available: strings.TrimSpace(line[availableStartIndex:sourceStartIndex]),
		})
	}

	return upgradeablePackages
}
