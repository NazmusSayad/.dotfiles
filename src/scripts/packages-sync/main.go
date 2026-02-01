package main

import (
	"dotfiles/src/helpers"
	"fmt"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println("⬇️ ", aurora.Faint("Installing Scoop packages..."))
	helpers.ExecNativeCommand([]string{"scoop-install"})
	fmt.Println()

	fmt.Println("⬇️ ", aurora.Faint("Installing Mise packages..."))
	helpers.ExecNativeCommand([]string{"mise", "install"})
	fmt.Println()

	fmt.Println("⬆️ ", aurora.Faint("Updating Scoop..."))
	helpers.ExecNativeCommand([]string{"scoop", "update"})
	fmt.Println()

	fmt.Println("⬆️ ", aurora.Faint("Updating Mise..."))
	helpers.ExecNativeCommand([]string{"mise", "upgrade"})
	fmt.Println()

	fmt.Println("⬆️ ", aurora.Faint("Updating Pacman..."))
	helpers.ExecNativeCommand([]string{"pacman", "-Syu", "--noconfirm"})
	fmt.Println()

	fmt.Println("🗑️ ", aurora.Faint("Pruning Scoop packages..."))
	helpers.ExecNativeCommand([]string{"scoop-prune"})
	fmt.Println()

	fmt.Println("🗑️ ", aurora.Faint("Pruning Mise packages..."))
	helpers.ExecNativeCommand([]string{"mise", "prune"})
	fmt.Println()

	fmt.Println("🗑️ ", aurora.Green("Pruning Pacman packages..."))
	helpers.ExecNativeCommand([]string{"pacman", "-Scc", "--noconfirm"})
	fmt.Println()
}
