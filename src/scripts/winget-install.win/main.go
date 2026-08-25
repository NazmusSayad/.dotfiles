package main

import (
	"fmt"
	"strconv"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/winget"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	packages := helpers.GetWingetApps()
	fmt.Println(aurora.Faint("Installing packages, total: " + strconv.Itoa(len(packages))))

	for _, p := range packages {
		if p.SkipInstall {
			fmt.Println()
			fmt.Println(aurora.Faint("- Skipping " + p.ID))
			continue
		}

		fmt.Println()
		fmt.Println(aurora.Faint("- Installing " + p.ID))

		helpers.ExecNativeCommand(
			append([]string{"winget", "install", "--no-upgrade"}, winget.BuildWingetOptions(p, p.InteractiveInstall)...),
			helpers.ExecCommandOptions{
				AsAdmin: p.ForceAdminInstall,
			},
		)
	}
}
