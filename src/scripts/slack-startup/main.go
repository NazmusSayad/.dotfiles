package main

import (
	slack "dotfiles/src/helpers/slack"
)

func main() {
	slack.SlackLaunch(slack.GetSlackStartupConfig())
}
