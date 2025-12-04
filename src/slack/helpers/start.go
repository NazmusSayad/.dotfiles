package slack_helpers

import (
	helpers "dotfiles/src"
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

	err = helpers.DetachedExec(runtimePath, "--startup")
	if err != nil {
		println("Error starting Slack")
	}
}
