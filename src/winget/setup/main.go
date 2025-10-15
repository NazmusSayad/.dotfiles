package main

import (
	helpers "dotfiles/src"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type entry struct {
	Name string
	Path string
}

func main() {
	helpers.EnsureAdminExecution()

	base, _ := os.Getwd()
	entries := []entry{
		{Name: "WinGet Install", Path: "___winget-install.exe"},
		{Name: "WinGet Upgrade", Path: "___winget-upgrade.exe"},
		{Name: "WinGet Auto Upgrade", Path: "___winget-upgrade-auto.exe"},
	}

	startMenu := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs")
	folder := filepath.Join(startMenu, "Winget")
	_ = os.MkdirAll(folder, 0755)

	for _, e := range entries {
		target := filepath.Join(base, e.Path)
		lnk := filepath.Join(folder, e.Name+".lnk")
		createShortcut(target, lnk)
	}
}

func createShortcut(target string, shortcut string) {
	psScript := []string{
		"$WshShell = New-Object -ComObject WScript.Shell",
		"$Shortcut = $WshShell.CreateShortcut(\"" + shortcut + "\")",
		"$Shortcut.WorkingDirectory = \"%USERPROFILE%\\Desktop\"",
		"$Shortcut.TargetPath = \"" + target + "\"",
		"$Shortcut.Save()",
	}

	cmd := exec.Command("powershell", "-c", strings.Join(psScript, ";"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}
