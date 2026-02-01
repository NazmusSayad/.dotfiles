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
	helpers.ExecNativeCommand([]string{"mise", "install"})
	fmt.Println()

	fmt.Println("▼", aurora.Faint("Installing Msys2 packages..."))
	msys2Packages := helpers.ReadConfig[[]string]("@/config/msys2-packages.jsonc")
	if len(msys2Packages) > 0 {
		fmt.Println(aurora.Faint("- Installing"), aurora.Green(strings.Join(msys2Packages, " ")))
		helpers.ExecNativeCommand(append([]string{"pacman", "--noconfirm", "-S", "--needed"}, msys2Packages...))
	}
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Scoop..."))
	helpers.ExecNativeCommand([]string{"scoop", "update"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Mise..."))
	helpers.ExecNativeCommand([]string{"mise", "upgrade"})
	fmt.Println()

	fmt.Println("△", aurora.Faint("Updating Pacman..."))
	helpers.ExecNativeCommand([]string{"pacman", "--noconfirm", "-Syu"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Scoop..."))
	scoop.PruneScoopApps()
	helpers.ExecNativeCommand([]string{"scoop", "cache", "rm", "*"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Mise..."))
	helpers.ExecNativeCommand([]string{"mise", "prune", "--yes"})
	helpers.ExecNativeCommand([]string{"mise", "cache", "clear", "--yes"})
	fmt.Println()

	fmt.Println("✘", aurora.Faint("Cleaning Pacman..."))
	helpers.ExecNativeCommand([]string{"pacman", "--noconfirm", "-Scc"})
	fmt.Println()
}
