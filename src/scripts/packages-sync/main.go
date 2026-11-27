package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/linuxpkg"

	"github.com/logrusorgru/aurora/v4"
)

var packageManagers = []string{"brew", "scoop", "pacman", "apt", "snap", "deb", "script", "mise"}

func main() {
	requested := os.Args[1:]
	for _, manager := range requested {
		if !slices.Contains(packageManagers, manager) {
			fmt.Println(aurora.Red("Unknown package manager: " + manager + " (expected: " + strings.Join(packageManagers, ", ") + ")"))
			os.Exit(1)
		}
	}

	selected := []string{}
	for _, manager := range packageManagers {
		if len(requested) > 0 && !slices.Contains(requested, manager) {
			continue
		}

		if !isAvailable(manager) {
			if len(requested) > 0 {
				fmt.Println("⚠", aurora.Yellow(manager+" is not available on this system, skipping."))
				fmt.Println()
			}
			continue
		}

		selected = append(selected, manager)
	}

	if slices.Contains(selected, "brew") {
		homebrewSync()
	}
	if slices.Contains(selected, "scoop") {
		scoopSync()
	}
	if slices.Contains(selected, "pacman") {
		pacmanSync()
	}
	if slices.Contains(selected, "apt") || slices.Contains(selected, "snap") || slices.Contains(selected, "deb") || slices.Contains(selected, "script") {
		linuxpkg.Sync(selected)
	}
	if slices.Contains(selected, "mise") {
		miseSync()
	}

	if runtime.GOOS == "windows" {
		wslSync(requested)
	}
}

func isAvailable(manager string) bool {
	switch manager {
	case "mise":
		return true
	case "brew":
		return runtime.GOOS == "darwin"
	case "scoop", "pacman":
		return runtime.GOOS == "windows"
	case "apt", "snap", "deb", "script":
		return runtime.GOOS == "linux" && !helpers.IsWSL()
	}

	panic("unknown package manager: " + manager)
}

func miseSync() {
	githubToken := helpers.GetGitHubToken()
	if githubToken == "" {
		fmt.Println("⚠", aurora.Yellow("No GitHub token found."))
		fmt.Println()
	}

	if runtime.GOOS != "windows" {
		fmt.Println("△", aurora.Faint("Updating Mise..."))
		runCommand(
			[]string{"mise", "self-update", "--yes"},
			helpers.ExecCommandOptions{ExtraEnv: map[string]string{"GITHUB_TOKEN": githubToken}},
		)
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
	runCommand([]string{"mise", "prune", "--yes"})

	fmt.Println("✘", aurora.Faint("Cleaning Mise cache..."))
	runCommand([]string{"mise", "cache", "clear", "--yes"})
}

func wslSync(requested []string) {
	fmt.Println(aurora.Blue("Syncing WSL packages"))
	cmd := exec.Command(
		"wsl.exe",
		"--distribution", "Ubuntu-26.04",
		"--user", "sayad",
		"--cd", "~",
		"--exec", "/bin/bash", "-lc",
		`cd "$HOME/.dotfiles" && "$HOME/.local/bin/mise" exec -- ./.local/bin.linux/psy `+strings.Join(requested, " "),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(aurora.Red("Failed to sync WSL packages: " + err.Error()))
	}
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

func scoopSync() {
	fmt.Println("◯", aurora.Faint("Configuring Scoop..."))
	runCommand([]string{"scoop", "config", "no_junction", "true"})

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
}

func pacmanSync() {
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
