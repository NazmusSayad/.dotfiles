package symlink

import (
	"slices"

	helpers "dotfiles/src/helpers"
)

type claudeAccountsConfig struct {
	ClaudeMd   string   `yaml:"CLAUDE.md"`
	ClaudeJson string   `yaml:"claude.json"`
	Accounts   []string `yaml:"accounts"`
}

var claudeSharedItems = map[string]string{
	"history.jsonl":   "File",
	"jobs":            "Folder",
	"plans":           "Folder",
	"tasks":           "Folder",
	"cache":           "Folder",
	"state":           "Folder",
	"plugins":         "Folder",
	"projects":        "Folder",
	"sessions":        "Folder",
	"paste-cache":     "Folder",
	"session-env":     "Folder",
	"shell-snapshots": "Folder",
	"file-history":    "Folder",
	"teams":           "Folder",
	"agents":          "Folder",
	"commands":        "Folder",
	"output-styles":   "Folder",
}

func readClaudeAccounts() map[string]Entry {
	config := helpers.ReadConfig[claudeAccountsConfig]("@/config/symlink.claude.yml")
	if config.ClaudeMd == "" || config.ClaudeJson == "" {
		panic("SYMLINK: CLAUDE.md or claude.json is missing in symlink.claude.yml")
	}

	primary := "~/.claude"
	var accounts []string
	for _, account := range config.Accounts {
		accounts = append(accounts, primary+"-"+account)
	}

	links := map[string]string{
		"~/.agents/skills": "skills",
		config.ClaudeMd:    "CLAUDE.md",
		config.ClaudeJson:  "settings.json",
	}

	entries := map[string]Entry{}
	for source, name := range links {
		entry := Entry{}
		for _, account := range slices.Concat([]string{primary}, accounts) {
			entry.Link = append(entry.Link, account+"/"+name)
		}
		entries[source] = entry
	}
	for name, create := range claudeSharedItems {
		entry := Entry{Create: create}
		for _, account := range accounts {
			entry.Link = append(entry.Link, account+"/"+name)
		}
		entries[primary+"/"+name] = entry
	}

	return entries
}
