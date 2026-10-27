package main

import (
	slack "dotfiles/src/helpers/slack"
)

func main() {
	config := slack.GetSlackStartupConfig()
	slack.SlackUpdatePresence(config)
}
