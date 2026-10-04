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

type agentConfig struct {
	Dir         string `yaml:"dir"`
	Model       string `yaml:"model"`
	MinDuration int    `yaml:"minDuration"`
}

type config struct {
	Codex  []agentConfig `yaml:"codex"`
	Claude []agentConfig `yaml:"claude"`
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

	for _, agent := range config.Codex {
		arguments := []string{"--no-daemon"}
		if agent.Model != "" {
			arguments = append(arguments, "--model", agent.Model)
		}
		arguments = append(arguments, "exec", "--skip-git-repo-check", "ping")
		run(lastRuns, logPath, "codex", "CODEX_HOME", agent, arguments...)
	}

	for _, agent := range config.Claude {
		arguments := []string{}
		if agent.Model != "" {
			arguments = append(arguments, "--model", agent.Model)
		}
		arguments = append(arguments, "-p", "ping")
		run(lastRuns, logPath, "claude", "CLAUDE_CONFIG_DIR", agent, arguments...)
	}
}

func run(lastRuns map[string]time.Time, logPath string, command string, environmentVariable string, agent agentConfig, arguments ...string) {
	resolvedProfile := helpers.ResolvePath(agent.Dir)
	account := command + ":" + resolvedProfile
	if lastRun, exists := lastRuns[account]; exists && time.Since(lastRun) < time.Duration(agent.MinDuration)*time.Minute {
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
		fmt.Printf("%s (%s): %v\n", command, agent.Dir, err)
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
