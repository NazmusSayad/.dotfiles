package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	slack "dotfiles/src/helpers/slack"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	data, err := os.ReadFile(filepath.Join(homeDir, ".slack-status"))
	if err != nil {
		fmt.Println("Error reading slack status file:", err)
		slack.SlackLaunch("work-hours")
		os.Exit(1)
	}

	status := strings.TrimSpace(string(data))
	slack.SlackLaunch(slack.SlackStatus(status))
}
