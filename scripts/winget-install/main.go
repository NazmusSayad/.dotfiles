package main

import (
	"dotfiles/src/helpers/winget"
	"fmt"
	"os"
	"os/exec"
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
		parts := winget.BuildWingetInstallCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	fmt.Println()
	fmt.Println(aurora.Green("Done!"))
}
