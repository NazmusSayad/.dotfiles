package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	slack_helpers "dotfiles/src/slack/helpers"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		println("Error getting user home directory:", err)
		time.Sleep(2000)
		os.Exit(1)
	}

	data, err := os.ReadFile(filepath.Join(homeDir, ".slack-status"))
	if err != nil {
		println("Error reading slack status file:", err)
		slack_helpers.SlackLaunch("work-hours")
		time.Sleep(2000)
		os.Exit(1)
	}

	status := strings.TrimSpace(string(data))
	slack_helpers.SlackLaunch(status)
}
