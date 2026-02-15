package main

import (
	"dotfiles/src/helpers"
	"fmt"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println("▼", aurora.Faint("Installing Scoop packages..."))
	runCommand([]string{"scoop-install"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Mise packages..."))
	runCommand([]string{"mise", "install"})
	fmt.Println()

	fmt.Println("◯", aurora.Faint("Preparing Pacman..."))
	runCommand([]string{"pacman", "--noconfirm", "-Sy"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Pacman packages..."))
	runCommand([]string{"msys-install"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Scoop..."))
	runCommand([]string{"scoop", "update"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Mise..."))
	runCommand([]string{"mise", "upgrade"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Pacman..."))
	runCommand([]string{"pacman", "--noconfirm", "-Su"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Scoop..."))
	runCommand([]string{"scoop-prune"})
	runCommand([]string{"scoop", "cache", "rm", "*"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Mise..."))
	runCommand([]string{"mise", "prune", "--yes"})
	runCommand([]string{"mise", "cache", "clear", "--yes"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Pacman..."))
	runCommand([]string{"pacman", "--noconfirm", "-Scc"})
	fmt.Println()
}

func runCommand(commands []string) {
	helpers.ExecNativeCommand(commands, helpers.ExecCommandOptions{Simulate: true})
}
