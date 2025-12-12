package main

import (
	"dotfiles/src/constants"
	"dotfiles/src/helpers"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	scriptsDir := helpers.ResolvePath(constants.SOURCE_DIR_SCRIPTS)
	startMenuDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "dotfiles")
	escape := func(s string) string { return strings.ReplaceAll(s, "'", "''") }

	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		_, err := exec.LookPath(entry.Name())
		if err != nil {
			continue
		}

		os.MkdirAll(startMenuDir, 0755)
		println("Installing", entry.Name())

		shortcutPath := filepath.Join(startMenuDir, entry.Name()+".lnk")
		targetCommand := "vbproxy user " + entryName

		cmd := exec.Command(
			"powershell",
			"-NoProfile",
			"-NonInteractive",
			"-Command",
			"$s='"+escape(targetCommand)+"';$t='"+escape(shortcutPath)+"';$ws=New-Object -ComObject WScript.Shell;$sc=$ws.CreateShortcut($t);$sc.TargetPath=$s;$sc.Save()",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
