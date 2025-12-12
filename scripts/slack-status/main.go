package main

import (
	"os"
	"path/filepath"
	"strings"

	slack_helpers "dotfiles/src/helpers/slack"

	"github.com/logrusorgru/aurora/v4"
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
		println("> " + label + ": " + aurora.Green("Always On").String())
	case slack_helpers.SlackStatusWorkTime:
		println("> " + label + ": " + aurora.Yellow("Work Time").String())
	case slack_helpers.SlackStatusDisabled:
		println("> " + label + ": " + aurora.Red("Disabled").String())
	}
}

func main() {
	initialStatus := readSlackStatus()
	renderSlackStatus("Current Slack Status", initialStatus)

	prompt := promptui.Select{
		Label:        "Select when to start Slack",
		Items:        []string{"Always", "Work Time", "Disabled"},
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Active: ("> " + aurora.BrightGreen("{{ . | green }}").String()),
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
}
