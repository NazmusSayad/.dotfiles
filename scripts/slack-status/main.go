package main

import (
	"os"
	"path/filepath"
	"strings"

	helpers "dotfiles/src/helpers"
	slack_helpers "dotfiles/src/helpers/slack"

	"github.com/manifoldco/promptui"
)

func getSlackStatusFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".slack-status")
}

func readSlackStatus() slack_helpers.SlackStatus {
	data, err := os.ReadFile(getSlackStatusFilePath())
	if err != nil {
		return slack_helpers.SlackStatusWorkTime
	}

	status := strings.TrimSpace(string(data))
	return slack_helpers.SlackStatus(status)
}

func writeSlackStatus(status slack_helpers.SlackStatus) {
	renderSlackStatus("Updating slack status to", status)
	os.WriteFile(getSlackStatusFilePath(), []byte(status), 0644)
	slack_helpers.SlackLaunch(status)
}

func renderSlackStatus(label string, status slack_helpers.SlackStatus) {

	switch status {
	case slack_helpers.SlackStatusAlways:
		println("> " + label + ": \033[32mAlways On\033[0m")
	case slack_helpers.SlackStatusWorkTime:
		println("> " + label + ": \033[33mWork Time\033[0m")
	case slack_helpers.SlackStatusDisabled:
		println("> " + label + ": \033[31mDisabled\033[0m")
	}
}

func main() {
	initialStatus := readSlackStatus()
	renderSlackStatus("Slack Status", initialStatus)

	prompt := promptui.Select{
		Label:        "Select when to start Slack",
		Items:        []string{"Always", "Work Time", "Disabled"},
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Active: "\033[32m> {{ . | green }}\033[0m",
		},
	}

	result, _, err := prompt.Run()
	if err != nil {
		return
	}

	switch result {
	case 0:
		writeSlackStatus(slack_helpers.SlackStatusAlways)
	case 1:
		writeSlackStatus(slack_helpers.SlackStatusWorkTime)
	case 2:
		writeSlackStatus(slack_helpers.SlackStatusDisabled)
	}

	helpers.PressAnyKeyOrWaitToExit()
}
