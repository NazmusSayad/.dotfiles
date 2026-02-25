package slack

import (
	helpers "dotfiles/src/helpers"
	"dotfiles/src/utils"
	"os"
)

const SlackStatusFileName = "~/.slack-startup-config"
const DefaultSlackStartupConfig = SlackStatusWorkTime

var resolvedPath = helpers.ResolvePath(SlackStatusFileName)

func GetSlackStartupConfig() SlackStatus {
	if !utils.IsFileExists(resolvedPath) {
		return DefaultSlackStartupConfig
	}

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return DefaultSlackStartupConfig
	}

	return SlackStatus(string(content))
}

func WriteSlackStartupConfig(config SlackStatus) {
	os.WriteFile(resolvedPath, []byte(config), 0644)
}
