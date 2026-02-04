package slack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	helpers "dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func IsSlackRunning() bool {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "Get-Process -Name 'slack' -ErrorAction SilentlyContinue")
	return cmd.Run() == nil
}

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

func SlackApplicationStart() {
	if IsSlackRunning() {
		fmt.Println(aurora.Faint("> Slack is already running"))
		return
	}

	runtimePath, err := GetSlackRuntimePath()
	if err != nil {
		fmt.Println("Error: Failed to get application runtime path")
		return
	}

	err = helpers.ExecNativeCommand(
		[]string{runtimePath, "--startup"},
		helpers.ExecCommandOptions{
			Detached: true,
		},
	)
	if err != nil {
		fmt.Println("Error: Failed to start Slack")
	}
}

func SlackApplicationStop() {
	if !IsSlackRunning() {
		fmt.Println(aurora.Faint("> Slack is not running"))
		return
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/IM", "slack.exe", "/F", "/T")
		cmd.Run()
	} else {
		cmd := exec.Command("pkill", "slack")
		cmd.Run()
	}
}
