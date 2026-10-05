package main

import (
	"encoding/json"
	"log"
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

type openCodeOpenAIConfig struct {
	Credential  string `yaml:"credential"`
	Model       string `yaml:"model"`
	MinDuration int    `yaml:"minDuration"`
}

type config struct {
	Codex          []agentConfig          `yaml:"codex"`
	OpenCodeOpenAI []openCodeOpenAIConfig `yaml:"opencode-openai"`
	Claude         []agentConfig          `yaml:"claude"`
}

func main() {
	config := helpers.ReadConfig[config]("@/config/coding-agents.yml")
	logDirectory := helpers.ResolvePath("~/.logs")
	if err := os.MkdirAll(logDirectory, 0o755); err != nil {
		panic(err)
	}
	statePath := filepath.Join(logDirectory, "agents-ping.json")
	logPath := filepath.Join(logDirectory, "agents-ping.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)
	logger.Println("agents ping started")
	defer logger.Println("agents ping finished")

	lastRuns := map[string]time.Time{}

	data, err := os.ReadFile(statePath)
	if err == nil {
		if err := json.Unmarshal(data, &lastRuns); err != nil {
			logger.Printf("%s: failed to decode ping state: %v", statePath, err)
			return
		}
	} else if !os.IsNotExist(err) {
		logger.Printf("%s: failed to read ping state: %v", statePath, err)
		return
	}

	for _, agent := range config.Codex {
		arguments := []string{"--no-daemon"}
		if agent.Model != "" {
			arguments = append(arguments, "--model", agent.Model)
		}
		arguments = append(arguments, "exec", "--skip-git-repo-check", "ping")
		run(lastRuns, statePath, logger, "codex", "CODEX_HOME", agent, arguments...)
	}

	for _, agent := range config.OpenCodeOpenAI {
		runOpenCodeOpenAI(lastRuns, statePath, logger, agent)
	}

	for _, agent := range config.Claude {
		arguments := []string{}
		if agent.Model != "" {
			arguments = append(arguments, "--model", agent.Model)
		}
		arguments = append(arguments, "-p", "ping")
		run(lastRuns, statePath, logger, "claude", "CLAUDE_CONFIG_DIR", agent, arguments...)
	}
}

func runOpenCodeOpenAI(lastRuns map[string]time.Time, statePath string, logger *log.Logger, agent openCodeOpenAIConfig) {
	if agent.Credential == "" {
		logger.Println("opencode-openai: credential is required")
		return
	}
	if agent.Model == "" {
		logger.Printf("opencode-openai (%s): model is required", agent.Credential)
		return
	}

	account := "opencode-openai:" + agent.Credential
	if lastRun, exists := lastRuns[account]; exists && time.Since(lastRun) < time.Duration(agent.MinDuration)*time.Minute {
		logger.Printf("opencode-openai (%s): skipped; last successful ping was %s", agent.Credential, lastRun.Format(time.RFC3339))
		return
	}
	logger.Printf("opencode-openai (%s, %s): starting", agent.Credential, agent.Model)

	directory, err := os.MkdirTemp("", "agents-ping-")
	if err != nil {
		logger.Printf("opencode-openai (%s): failed to create temporary directory: %v", agent.Credential, err)
		return
	}
	defer os.RemoveAll(directory)

	switchCommand := exec.Command("opencode", "auth", "switch", "openai", agent.Credential)
	switchCommand.Dir = directory
	switchCommand.Stdout = logger.Writer()
	switchCommand.Stderr = logger.Writer()
	if err := switchCommand.Run(); err != nil {
		logger.Printf("opencode-openai (%s): account switch failed: %v", agent.Credential, err)
		return
	}
	logger.Printf("opencode-openai (%s): account switched", agent.Credential)

	pingCommand := exec.Command("opencode", "run", "--model", "openai/"+agent.Model, "ping")
	pingCommand.Dir = directory
	pingCommand.Stdout = logger.Writer()
	pingCommand.Stderr = logger.Writer()
	if err := pingCommand.Run(); err != nil {
		logger.Printf("opencode-openai (%s): ping failed: %v", agent.Credential, err)
		return
	}

	logger.Printf("opencode-openai (%s): ping succeeded", agent.Credential)
	recordRun(lastRuns, statePath, logger, account)
}

func run(lastRuns map[string]time.Time, statePath string, logger *log.Logger, command string, environmentVariable string, agent agentConfig, arguments ...string) {
	resolvedProfile := helpers.ResolvePath(agent.Dir)
	account := command + ":" + resolvedProfile
	if lastRun, exists := lastRuns[account]; exists && time.Since(lastRun) < time.Duration(agent.MinDuration)*time.Minute {
		logger.Printf("%s (%s): skipped; last successful ping was %s", command, agent.Dir, lastRun.Format(time.RFC3339))
		return
	}
	logger.Printf("%s (%s, %s): starting", command, agent.Dir, agent.Model)

	directory, err := os.MkdirTemp("", "agents-ping-")
	if err != nil {
		logger.Printf("%s (%s): failed to create temporary directory: %v", command, agent.Dir, err)
		return
	}
	defer os.RemoveAll(directory)

	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory
	cmd.Env = append(os.Environ(), environmentVariable+"="+resolvedProfile)
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()

	if err := cmd.Run(); err != nil {
		logger.Printf("%s (%s): ping failed: %v", command, agent.Dir, err)
		return
	}

	logger.Printf("%s (%s): ping succeeded", command, agent.Dir)
	recordRun(lastRuns, statePath, logger, account)
}

func recordRun(lastRuns map[string]time.Time, statePath string, logger *log.Logger, account string) {
	lastRuns[account] = time.Now()
	data, err := json.MarshalIndent(lastRuns, "", "  ")
	if err != nil {
		logger.Printf("%s: failed to encode ping state: %v", statePath, err)
		return
	}
	if err := os.WriteFile(statePath, data, 0o644); err != nil {
		logger.Printf("%s: failed to write ping state: %v", statePath, err)
	}
}
