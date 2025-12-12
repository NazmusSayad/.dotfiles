package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/helpers/winget"
	"fmt"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	packages := winget.GetWingetPackages()
	fmt.Println(aurora.Faint("Installing packages, total: " + strconv.Itoa(len(packages))))

	for _, p := range packages {
		if p.SkipInstall {
			fmt.Println()
			fmt.Println(aurora.Faint("- Skipping " + p.ID))
			continue
		}

		fmt.Println()
		fmt.Println(aurora.Faint("- Installing " + p.ID))

		args := winget.BuildWingetInstallCommands(p)
		helpers.ExecNativeCommand(helpers.ExecCommandOptions{
			Command: args[0],
			Args:    args[1:],
		})
	}
}
