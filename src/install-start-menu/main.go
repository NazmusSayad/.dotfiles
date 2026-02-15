package main

import (
	"dotfiles/src/constants"
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func escape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func main() {
	proxyPausePath, err := exec.LookPath("proxy-pause")
	if err != nil {
		fmt.Println("proxy-pause not found in PATH")
		os.Exit(1)
	}

	startMenuDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "dotfiles")
	if utils.IsFileExists(startMenuDir) {
		fmt.Println("Removing", startMenuDir)
		os.RemoveAll(startMenuDir)
	}

	os.MkdirAll(startMenuDir, 0755)

	for scriptName, script := range constants.BIN_SCRIPTS {
		if script.StartMenu == "" {
			continue
		}

		fmt.Println("> Installing", scriptName, script.StartMenu)

		targetCommand := `"` + proxyPausePath + `"`
		shortcutPath := filepath.Join(startMenuDir, script.StartMenu+".lnk")
		arguments := `"` + utils.Ternary(script.Exe != "", script.Exe, scriptName) + `"`

		helpers.ExecNativeCommand([]string{
			"powershell",
			"-NoProfile",
			"-NonInteractive",
			"-Command",
			"$s='" + escape(targetCommand) + "';$a='" + escape(arguments) + "';$t='" + escape(shortcutPath) + "';$ws=New-Object -ComObject WScript.Shell;$sc=$ws.CreateShortcut($t);$sc.TargetPath=$s;$sc.Arguments=$a;$sc.Save()",
		})
	}
}
