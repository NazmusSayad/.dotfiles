package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	slack "dotfiles/src/helpers/slack"

	"github.com/logrusorgru/aurora/v4"
	"github.com/manifoldco/promptui"
)

func getSlackStatusFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".slack-status")
}

func readSlackStatus() slack.SlackStatus {
	data, err := os.ReadFile(getSlackStatusFilePath())
	if err != nil {
		return slack.SlackStatusWorkTime
	}

	status := strings.TrimSpace(string(data))
	return slack.SlackStatus(status)
}

func writeSlackStatus(status slack.SlackStatus) {
	renderSlackStatus("Updating slack status to", status)
	os.WriteFile(getSlackStatusFilePath(), []byte(status), 0644)
	slack.SlackLaunch(status)
}

func renderSlackStatus(label string, status slack.SlackStatus) {

	switch status {
	case slack.SlackStatusAlways:
		fmt.Println("> " + label + ": " + aurora.Green("Always On").String())
	case slack.SlackStatusWorkTime:
		fmt.Println("> " + label + ": " + aurora.Yellow("Work Time").String())
	case slack.SlackStatusDisabled:
		fmt.Println("> " + label + ": " + aurora.Red("Disabled").String())
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
			Active: "> " + aurora.BrightGreen("{{ . | green }}").String(),
		},
	}

	result, _, err := prompt.Run()
	if err != nil {
		return
	}

	switch result {
	case 0:
		writeSlackStatus(slack.SlackStatusAlways)
	case 1:
		writeSlackStatus(slack.SlackStatusWorkTime)
	case 2:
		writeSlackStatus(slack.SlackStatusDisabled)
	}
}
