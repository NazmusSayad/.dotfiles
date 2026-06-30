package main

import (
	"fmt"
	"runtime"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if runtime.GOOS == "darwin" {
		macosSync()
	} else if runtime.GOOS == "windows" {
		windowsSync()
	}
}

func macosSync() {
	// Brew
	brewfilePath := helpers.ResolvePath("@/config/Brewfile")
	brewFileTaps := helpers.GetBrewTaps(brewfilePath)
	if len(brewFileTaps) > 0 {
		fmt.Println("◯", aurora.Faint("Trusting Brew taps..."))

		for _, tap := range brewFileTaps {
			helpers.ExecNativeCommand([]string{"brew", "trust", tap})
		}

		fmt.Println()
	}

	fmt.Println("△", aurora.Faint("Updating Brew..."))
	helpers.ExecNativeCommand([]string{"brew", "update"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Brew packages..."))
	helpers.ExecNativeCommand([]string{"brew", "bundle", "install", "--file=" + brewfilePath})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Uninstalling Brew packages..."))
	helpers.ExecNativeCommand([]string{"brew", "bundle", "cleanup", "--force", "--file=" + brewfilePath})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Upgrading Brew Apps..."))
	helpers.ExecNativeCommand([]string{"brew", "upgrade", "--yes"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Brew..."))
	helpers.ExecNativeCommand([]string{"brew", "cleanup", "--prune=all", "-s"})
	fmt.Println()

	// Mise
	fmt.Println("▼", aurora.Faint("Installing Mise packages..."))
	helpers.ExecNativeCommand([]string{"mise", "install"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Mise..."))
	helpers.ExecNativeCommand([]string{"mise", "upgrade"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Mise..."))
	helpers.ExecNativeCommand([]string{"mise", "prune", "--yes"})
	helpers.ExecNativeCommand([]string{"mise", "cache", "clear", "--yes"})
	fmt.Println()
}

func windowsSync() {
	// Scoop
	fmt.Println("✘", aurora.Faint("Uninstalling Scoop Apps..."))
	runCommand([]string{"scoop-prune"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Scoop..."))
	runCommand([]string{"scoop", "update", "--quiet"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Scoop packages..."))
	runCommand([]string{"scoop-install"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Scoop Apps..."))
	runCommand([]string{"scoop", "update", "--all", "--quiet"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Scoop..."))
	runCommand([]string{"scoop", "cache", "rm", "*"})
	fmt.Println()

	// Pacman
	fmt.Println("◯", aurora.Faint("Preparing Pacman..."))
	runCommand([]string{"pacman", "-Sy", "--noconfirm"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Pacman packages..."))
	runCommand([]string{"msys-install"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Pacman..."))
	runCommand([]string{"pacman", "-Su", "--noconfirm"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Pacman..."))
	runCommand([]string{"pacman", "-Scc", "--noconfirm"})
	fmt.Println()

	// Mise
	fmt.Println("▼", aurora.Faint("Installing Mise packages..."))
	runCommand([]string{"mise", "install"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Mise..."))
	runCommand([]string{"mise", "upgrade"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Mise..."))
	runCommand([]string{"mise", "prune", "--yes"})
	runCommand([]string{"mise", "cache", "clear", "--yes"})
	fmt.Println()
}

func runCommand(commands []string) {
	helpers.ExecNativeCommand(commands, helpers.ExecCommandOptions{Simulate: true})
}
