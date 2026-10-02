package main

import (
	"fmt"
	"os"
	"os/exec"

	"dotfiles/src/helpers"
)

type config struct {
	Codex  []string `yaml:"codex"`
	Claude []string `yaml:"claude"`
}

func main() {
	config := helpers.ReadConfig[config]("@/config/coding-agents.yml")

	for _, profile := range config.Codex {
		run("codex", "CODEX_HOME", profile, "--no-daemon", "exec", "--skip-git-repo-check", "hi")
	}

	for _, profile := range config.Claude {
		run("claude", "CLAUDE_CONFIG_DIR", profile, "-p", "hi")
	}
}

func run(command string, environmentVariable string, profile string, arguments ...string) {
	directory, err := os.MkdirTemp("", "agents-ping-")
	if err != nil {
		fmt.Printf("%s: %v\n", command, err)
		return
	}
	defer os.RemoveAll(directory)

	cmd := exec.Command(command, arguments...)
	cmd.Dir = directory
	cmd.Env = append(os.Environ(), environmentVariable+"="+helpers.ResolvePath(profile))

	if err := cmd.Run(); err != nil {
		fmt.Printf("%s (%s): %v\n", command, profile, err)
	}
}
