package slack_helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetSlackRuntimePath() (string, error) {
	slackPath := filepath.Join(os.Getenv("LOCALAPPDATA"), "slack")
	slackBaseExe := filepath.Join(slackPath, "slack.exe")

	out, err := exec.Command("powershell", "-NoProfile", "-Command", "(Get-Item '"+slackBaseExe+"').VersionInfo.ProductVersion").Output()
	if err != nil {
		return "", err
	}

	productVersion := strings.TrimSpace(string(out))
	runtimePath := filepath.Join(slackPath, "app-"+productVersion, "slack.exe")
	if _, err := os.Stat(runtimePath); err != nil {
		return "", err
	}

	return runtimePath, nil
}

func IsSlackRunning() bool {
	err := exec.Command("powershell", "-NoProfile", "-Command", "Get-Process -Name 'slack' -ErrorAction SilentlyContinue").Run()
	return err == nil
}
