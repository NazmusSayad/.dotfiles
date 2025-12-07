package slack_helpers

import (
	"os/exec"
	"runtime"
)

func SlackApplicationStop() {
	if !IsSlackRunning() {
		println("Slack is not running")
		return
	}

	println("Slack Stop")
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/IM", "slack.exe", "/F", "/T")
		cmd.Run()
	} else {
		cmd := exec.Command("pkill", "slack")
		cmd.Run()
	}
}
