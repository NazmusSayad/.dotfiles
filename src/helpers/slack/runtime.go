package slack

import (
	"fmt"
	"os/exec"
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
	scoopSlackCmd := exec.Command("scoop", "which", "slack")
	slackPathOutput, err := scoopSlackCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	slackPath := strings.TrimSpace(string(slackPathOutput))
	return helpers.ResolvePath(slackPath), nil
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
			NoWait:   true,
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
