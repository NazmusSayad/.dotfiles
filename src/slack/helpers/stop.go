package slack_helpers

import "os/exec"

func SlackApplicationStop() {
	if !IsSlackRunning() {
		println("Slack is not running")
		return
	}

	println("Slack Stop")
	exec.Command("powershell", "-NoProfile", "-Command", "Stop-Process -Name 'slack' -Force").Run()
}
