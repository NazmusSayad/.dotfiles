package main

import (
	"dotfiles/src/constants"
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	scriptsDir := helpers.ResolvePath(constants.SCRIPTS_SOURCE_DIR)
	startMenuDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "dotfiles")
	escape := func(s string) string { return strings.ReplaceAll(s, "'", "''") }

	proxyPausePath, err := exec.LookPath("proxy-pause")
	if err != nil {
		fmt.Println("proxy-pause not found in PATH")
		os.Exit(1)
	}

	if helpers.IsFileExists(startMenuDir) {
		fmt.Println("Removing", startMenuDir)
		os.RemoveAll(startMenuDir)
	}
	os.MkdirAll(startMenuDir, 0755)

	entries, err := os.ReadDir(scriptsDir)
	if err != nil {
		os.Exit(1)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryName := entry.Name()
		binPath, err := exec.LookPath(entryName)
		if err != nil {
			continue
		}

		fmt.Println("Installing", entryName)

		shortcutPath := filepath.Join(startMenuDir, entryName+".lnk")
		targetCommand := `"` + proxyPausePath + `"`
		arguments := `"` + binPath + `"`

		cmd := exec.Command(
			"powershell",
			"-NoProfile",
			"-NonInteractive",
			"-Command",
			"$s='"+escape(targetCommand)+"';$a='"+escape(arguments)+"';$t='"+escape(shortcutPath)+"';$ws=New-Object -ComObject WScript.Shell;$sc=$ws.CreateShortcut($t);$sc.TargetPath=$s;$sc.Arguments=$a;$sc.Save()",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
