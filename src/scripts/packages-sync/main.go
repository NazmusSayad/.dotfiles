package main

import (
	"fmt"
	"runtime"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	switch runtime.GOOS {
	case "darwin":
		homebrewSync()
	case "windows":
		windowsSync()
	}

	githubToken := helpers.GetGitHubToken()
	if githubToken == "" {
		fmt.Println("⚠", aurora.Yellow("No GitHub token found."))
		fmt.Println()
	}

	// Mise
	if runtime.GOOS == "darwin" {
		fmt.Println("△", aurora.Faint("Updating Mise..."))
		runCommand([]string{"mise", "self-update", "--yes"})
	}

	fmt.Println("✘", aurora.Faint("Uninstalling Mise packages..."))
	runCommand([]string{"mise", "prune", "--yes"})

	fmt.Println("▼", aurora.Faint("Installing Mise packages..."))
	runCommand(
		[]string{"mise", "install", "--yes"},
		helpers.ExecCommandOptions{ExtraEnv: map[string]string{"GITHUB_TOKEN": githubToken}},
	)

	fmt.Println("△", aurora.Faint("Updating Mise packages..."))
	runCommand(
		[]string{"mise", "upgrade", "--yes"},
		helpers.ExecCommandOptions{ExtraEnv: map[string]string{"GITHUB_TOKEN": githubToken}},
	)

	fmt.Println("✘", aurora.Faint("Cleaning Mise packages..."))
	runCommand([]string{"mise", "cache", "clear", "--yes"})
}

func homebrewSync() {
	brewfilePath := helpers.ResolvePath("@/config/Brewfile")
	brewFileTaps := helpers.GetBrewTaps(brewfilePath)
	if len(brewFileTaps) > 0 {
		fmt.Println("◯", aurora.Faint("Trusting Brew taps..."))
		for _, tap := range brewFileTaps {
			helpers.ExecNativeCommand([]string{"brew", "trust", tap})
		}
		fmt.Println()
	}

	fmt.Println("✘", aurora.Faint("Uninstalling Brew packages..."))
	runCommand([]string{"brew", "bundle", "cleanup", "--force", "--file=" + brewfilePath})

	fmt.Println("✘", aurora.Faint("Autoremove Brew dependencies..."))
	runCommand([]string{"brew", "autoremove"})

	if len(brewFileTaps) > 0 {
		fmt.Println("◯", aurora.Faint("Trusting Brew taps..."))
		for _, tap := range brewFileTaps {
			helpers.ExecNativeCommand([]string{"brew", "trust", tap})
		}
		fmt.Println()
	}

	fmt.Println("△", aurora.Faint("Updating Brew..."))
	runCommand([]string{"brew", "update"})

	fmt.Println("▼", aurora.Faint("Installing Brew packages..."))
	runCommand([]string{"brew", "bundle", "install", "--file=" + brewfilePath})

	fmt.Println("△", aurora.Faint("Upgrading Brew packages..."))
	runCommand([]string{"brew", "upgrade", "--yes"})

	fmt.Println("✘", aurora.Faint("Cleaning Brew..."))
	runCommand([]string{"brew", "cleanup", "--scrub"})
}

func windowsSync() {
	// Scoop
	fmt.Println("✘", aurora.Faint("Uninstalling Scoop Apps..."))
	runCommand([]string{"scoop-prune"})

	fmt.Println("△", aurora.Faint("Updating Scoop..."))
	runCommand([]string{"scoop", "update", "--quiet"})

	fmt.Println("▼", aurora.Faint("Installing Scoop packages..."))
	runCommand([]string{"scoop-install"})

	fmt.Println("△", aurora.Faint("Updating Scoop Apps..."))
	runCommand([]string{"scoop", "update", "--all", "--quiet"})

	fmt.Println("✘", aurora.Faint("Cleaning Scoop..."))
	runCommand([]string{"scoop", "cache", "rm", "*"})

	// Pacman
	fmt.Println("◯", aurora.Faint("Preparing Pacman..."))
	runCommand([]string{"pacman", "-Sy", "--noconfirm"})

	fmt.Println("▼", aurora.Faint("Installing Pacman packages..."))
	runCommand([]string{"msys-install"})

	fmt.Println("△", aurora.Faint("Updating Pacman..."))
	runCommand([]string{"pacman", "-Su", "--noconfirm"})

	fmt.Println("✘", aurora.Faint("Cleaning Pacman..."))
	runCommand([]string{"pacman", "-Scc", "--noconfirm"})
}

func runCommand(commands []string, options ...helpers.ExecCommandOptions) {
	opts := helpers.ExecCommandOptions{Simulate: true}
	if len(options) > 0 {
		opts = options[0]
		opts.Simulate = true
	}

	helpers.ExecNativeCommand(commands, opts)
	fmt.Println()
}
