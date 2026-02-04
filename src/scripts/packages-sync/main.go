package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/helpers/scoop"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println("▼", aurora.Faint("Installing Scoop packages..."))
	scoop.InstallScoopApps()
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Mise packages..."))
	runCommand([]string{"mise", "install"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Msys2 packages..."))
	msys2Packages := helpers.ReadConfig[[]string]("@/config/msys2-packages.jsonc")
	if len(msys2Packages) > 0 {
		fmt.Println(aurora.Faint("- Installing"), aurora.Green(strings.Join(msys2Packages, " ")))
		runCommand(append([]string{"pacman", "--noconfirm", "-S", "--needed"}, msys2Packages...))
	}
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Scoop..."))
	runCommand([]string{"scoop", "update"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Mise..."))
	runCommand([]string{"mise", "upgrade"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Pacman..."))
	runCommand([]string{"pacman", "--noconfirm", "-Syu"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Scoop..."))
	scoop.PruneScoopApps()
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
	helpers.ExecNativeCommand(commands, helpers.ExecCommandOptions{Detached: true})
}
