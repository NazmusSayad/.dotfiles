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

var escape = func(s string) string { return strings.ReplaceAll(s, "'", "''") }

func main() {

	proxyPausePath, err := exec.LookPath("proxy-pause")
	if err != nil {
		fmt.Println("proxy-pause not found in PATH")
		os.Exit(1)
	}

	startMenuDir := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "dotfiles")
	if helpers.IsFileExists(startMenuDir) {
		fmt.Println("Removing", startMenuDir)
		os.RemoveAll(startMenuDir)
	}
	os.MkdirAll(startMenuDir, 0755)

	for scriptName, script := range constants.SCRIPTS_MAP {
		if script.StartMenuName == "" {
			continue
		}

		fmt.Println("> Installing", scriptName, script.StartMenuName)

		shortcutPath := filepath.Join(startMenuDir, script.StartMenuName+".lnk")
		targetCommand := `"` + proxyPausePath + `"`
		arguments := `"` + scriptName + `"`

		helpers.ExecNativeCommand(helpers.ExecCommandOptions{
			Command: "powershell",
			Args: []string{
				"-NoProfile",
				"-NonInteractive",
				"-Command",
				"$s='" + escape(targetCommand) + "';$a='" + escape(arguments) + "';$t='" + escape(shortcutPath) + "';$ws=New-Object -ComObject WScript.Shell;$sc=$ws.CreateShortcut($t);$sc.TargetPath=$s;$sc.Arguments=$a;$sc.Save()",
			},
		})
	}
}
