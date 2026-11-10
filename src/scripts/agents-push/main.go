package main

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/src/helpers"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

var claudeItems = []string{
	"projects",
	"history.jsonl",
	"file-history",
	"plans",
	"tasks",
	"paste-cache",
}

func main() {
	if !utils.IsCommandInPath("restic") {
		fmt.Println(aurora.Red("Error: restic not installed"))
		os.Exit(1)
	}

	host := helpers.DetectClaudeDevice()
	repository := helpers.ClaudeStorage(host)

	claudeDir := helpers.ResolvePath("~/.claude")

	var items []string
	for _, item := range claudeItems {
		if utils.IsFileExists(filepath.Join(claudeDir, item)) {
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		fmt.Println(aurora.Red("Error: No Claude data found in " + claudeDir))
		os.Exit(1)
	}

	helpers.ExecNativeCommand(
		append([]string{"restic", "backup", "--repo", repository, "--insecure-no-password", "--host", host, "--tag", "claude"}, items...),
		helpers.ExecCommandOptions{Dir: claudeDir, Exit: true},
	)
}
