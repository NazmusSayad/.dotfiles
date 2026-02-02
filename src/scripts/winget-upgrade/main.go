package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/helpers/winget"
	"fmt"
	"slices"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	packages := helpers.ReadConfig[[]winget.WingetPackage]("@/config/winget-apps.jsonc")

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

		helpers.ExecNativeCommand(append([]string{"winget"}, winget.BuildWingetUpgradeArguments(p)...))
	}
}
