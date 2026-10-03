package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"dotfiles/src/helpers"
)

type config struct {
	Codex  []string `yaml:"codex"`
	Claude []string `yaml:"claude"`
}

func main() {
	config := helpers.ReadConfig[config]("@/config/coding-agents.yml")
	logDirectory := helpers.ResolvePath("~/.logs")
	if err := os.MkdirAll(logDirectory, 0o755); err != nil {
		panic(err)
	}
	logPath := filepath.Join(logDirectory, "agents-ping.json")
	lastRuns := map[string]time.Time{}

	data, err := os.ReadFile(logPath)
	if err == nil {
		if err := json.Unmarshal(data, &lastRuns); err != nil {
			panic(err)
		}
	} else if !os.IsNotExist(err) {
		panic(err)
	}

	for _, profile := range config.Codex {
		run(lastRuns, logPath, "codex", "CODEX_HOME", profile, "--no-daemon", "exec", "--skip-git-repo-check", "hi")
	}

	for _, profile := range config.Claude {
		run(lastRuns, logPath, "claude", "CLAUDE_CONFIG_DIR", profile, "-p", "hi")
	}
}

func run(lastRuns map[string]time.Time, logPath string, command string, environmentVariable string, profile string, arguments ...string) {
	resolvedProfile := helpers.ResolvePath(profile)
	account := command + ":" + resolvedProfile
	if lastRun, exists := lastRuns[account]; exists && time.Since(lastRun) < 30*time.Minute {
		return
	}

	directory, err := os.MkdirTemp("", "agents-ping-")
	if err != nil {
		fmt.Printf("%s: %v\n", command, err)
		return
	}
	defer os.RemoveAll(directory)

	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory
	cmd.Env = append(os.Environ(), environmentVariable+"="+resolvedProfile)

	if err := cmd.Run(); err != nil {
		fmt.Printf("%s (%s): %v\n", command, profile, err)
		return
	}

	lastRuns[account] = time.Now()
	data, err := json.MarshalIndent(lastRuns, "", "  ")
	if err != nil {
		fmt.Printf("%s: %v\n", logPath, err)
		return
	}
	if err := os.WriteFile(logPath, data, 0o644); err != nil {
		fmt.Printf("%s: %v\n", logPath, err)
	}
}
