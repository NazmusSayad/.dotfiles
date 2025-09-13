package winget

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type WingetPackage struct {
	ID      string
	Name    string
	Version string
}

func GetWingetPackages(path string) []WingetPackage {
	f, err := os.Open(path)
	if err != nil {
		return []WingetPackage{}
	}

	defer f.Close()

	var pkgs []WingetPackage
	dec := json.NewDecoder(f)
	if err := dec.Decode(&pkgs); err != nil {
		return []WingetPackage{}
	}

	return pkgs
}

func ConfirmIsAdminExec() {
	psCmd := `if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	if err := cmd.Run(); err == nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "This program requires administrator privileges.")
		fmt.Fprintln(os.Stderr, "Trying to relaunch with elevated privileges...")

		exePath, e := os.Executable()
		if e != nil {
			_, _ = reader.ReadString('\n')
			os.Exit(1)
		}

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+exePath+"' -Verb RunAs")
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to relaunch with elevated privileges. Press Enter to exit...")
			_, _ = reader.ReadString('\n')
			os.Exit(1)
		}

		os.Exit(0)
	}
}

func BuildWingetInstallCommands(p WingetPackage) []string {
	parts := []string{"winget", "install", p.ID}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	return append(parts, "--interactive", "--accept-package-agreements", "--accept-source-agreements", "--no-upgrade")
}

func BuildWingetUpgradeCommands(p WingetPackage) []string {
	parts := []string{"winget", "upgrade", p.ID}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	return append(parts, "--interactive", "--accept-package-agreements", "--accept-source-agreements", "--uninstall-previous")
}
