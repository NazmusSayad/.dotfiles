package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	slack "dotfiles/src/helpers/slack"
	"dotfiles/src/utils"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	statusFilePath := filepath.Join(homeDir, ".slack-status")
	if !utils.IsFileExists(statusFilePath) {
		slack.SlackLaunch(slack.SlackStatusWorkTime)
		os.Exit(0)
	}

	data, err := os.ReadFile(statusFilePath)
	if err != nil {
		fmt.Println("Error reading slack status file:", err)
		os.Exit(1)
	}

	status := strings.TrimSpace(string(data))
	slack.SlackLaunch(slack.SlackStatus(status))
}
