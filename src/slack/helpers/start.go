package slack_helpers

import (
	"os"
	"os/exec"
)

func SlackApplicationStart() {
	if IsSlackRunning() {
		println("Slack is already running")
		return
	}

	runtimePath, err := GetSlackRuntimePath()
	if err != nil {
		println("Error getting application runtime path")
		return
	}

	println("Slack Start", runtimePath)

	cmd := exec.Command(runtimePath, "--startup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}
