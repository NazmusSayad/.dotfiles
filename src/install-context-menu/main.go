package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"path"
)

func cleanupContextMenu(folder string) {
	helpers.ExecNativeCommand([]string{"reg", "delete", `HKEY_CURRENT_USER\Software\Classes\*\shell\	` + folder, "/f"})
	helpers.ExecNativeCommand([]string{"reg", "delete", `HKEY_CURRENT_USER\Software\Classes\Directory\shell\` + folder, "/f"})
	helpers.ExecNativeCommand([]string{"reg", "delete", `HKEY_CURRENT_USER\Software\Classes\Directory\Background\shell\` + folder, "/f"})
}

func installContextMenu(folder string, label string, cmd string) {
	regCode := `Windows Registry Editor Version 5.00

[HKEY_CURRENT_USER\Software\Classes\*\shell\` + folder + `]
@="` + label + `"
"Icon"="C:\\Users\\Sayad\\scoop\\apps\\` + folder + `\\current\\` + cmd + `"
[HKEY_CURRENT_USER\Software\Classes\*\shell\` + folder + `\command]
@="\"C:\\Users\\Sayad\\scoop\\apps\\` + folder + `\\current\\` + cmd + `\" \"%1\""

[HKEY_CURRENT_USER\Software\Classes\Directory\shell\` + folder + `]
@="` + label + `"
"Icon"="C:\\Users\\Sayad\\scoop\\apps\\` + folder + `\\current\\` + cmd + `"
[HKEY_CURRENT_USER\Software\Classes\Directory\shell\` + folder + `\command]
@="\"C:\\Users\\Sayad\\scoop\\apps\\` + folder + `\\current\\` + cmd + `\" \"%V\""

[HKEY_CURRENT_USER\Software\Classes\Directory\Background\shell\` + folder + `]
@="` + label + `"
"Icon"="C:\\Users\\Sayad\\scoop\\apps\\` + folder + `\\current\\` + cmd + `"
[HKEY_CURRENT_USER\Software\Classes\Directory\Background\shell\` + folder + `\command]
@="\"C:\\Users\\Sayad\\scoop\\apps\\` + folder + `\\current\\` + cmd + `\" \"%V\""`

	tempPath := path.Join(os.TempDir(), "./install-something.reg")
	os.WriteFile(tempPath, []byte(regCode), 0644)

	fmt.Println("> Installing context menu for", label, tempPath)
	helpers.ExecNativeCommand([]string{"reg", "import", tempPath})
}

func main() {
	if utils.IsCommandInPath("code") {
		installContextMenu("vscode", "Open in Code", "Code.exe")
	} else {
		cleanupContextMenu("vscode")
	}

	if utils.IsCommandInPath("cursor") {
		installContextMenu("cursor", "Open in Cursor", "Cursor.exe")
	} else {
		cleanupContextMenu("cursor")
	}
}
