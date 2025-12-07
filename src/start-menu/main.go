package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	helpers "dotfiles/src"
)

type AppConfig struct {
	Name string
	Path string
}

func execPsCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("powershell command failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func main() {
	helpers.EnsureAdminExecution()

	data, err := helpers.ReadDotfilesConfigJSONC("./config/start-menu-apps.jsonc")
	if err != nil {
		println("Failed to read config: %v\n", err)
		os.Exit(1)
	}

	var apps []AppConfig
	if err := json.Unmarshal(data, &apps); err != nil {
		println("Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	dest := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "dotfiles")

	if err := os.RemoveAll(dest); err != nil && !os.IsNotExist(err) {
		println("Failed to remove existing directory: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		println("Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	for _, app := range apps {
		println("Creating shortcut: %s.lnk -> %s\n", app.Name, app.Path)

		shortcutPath := filepath.Join(dest, app.Name+".lnk")
		targetPath := helpers.ResolvePath(app.Path)

		psCmd := fmt.Sprintf(`$w=New-Object -ComObject WScript.Shell; $s=$w.CreateShortcut('%s'); $s.TargetPath='%s'; $s.Save()`,
			strings.ReplaceAll(shortcutPath, "'", "''"),
			strings.ReplaceAll(targetPath, "'", "''"))

		_, err := execPsCommand(psCmd)
		if err != nil {
			println("FAIL: %s.lnk\n", app.Name)
		} else {
			println("OK: %s.lnk\n", app.Name)
		}
	}

	helpers.PressAnyKeyOrWaitToExit()
}
