package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/helpers/winget"
	"fmt"
	"slices"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	packages := winget.GetWingetPackages()

	var packageIDs []string
	upgradeablePackages := winget.GetUpgradeablePackages()

	for _, p := range upgradeablePackages {
		fmt.Println()
		fmt.Println("ID:", p.ID)
		fmt.Println("Current Version: " + aurora.Red(p.Version).String())
		fmt.Println("Available Version: " + aurora.Green(p.Available).String())
		packageIDs = append(packageIDs, p.ID)
	}

	for _, p := range packages {
		if !slices.Contains(packageIDs, p.ID) {
			continue
		}

		if p.SkipUpgrade || p.Version != "" {
			fmt.Println()
			fmt.Println(aurora.Faint("- Skipping " + p.ID))
			continue
		}

		fmt.Println()
		fmt.Println(aurora.Faint("- Upgrading " + p.ID))

		args := winget.BuildWingetUpgradeCommands(p)
		helpers.ExecNativeCommand(helpers.ExecCommandOptions{
			Command: args[0],
			Args:    args[1:],
		})
	}
}
