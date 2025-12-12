package main

import (
	helpers "dotfiles/src/helpers"
	"dotfiles/src/helpers/winget"
	"os"
	"os/exec"
)

func main() {
	packages := winget.GetWingetPackages()
	println("Installing packages, total:", len(packages))

	for _, p := range packages {
		if p.SkipInstall {
			println("\n- Skipping", p.ID)
			continue
		}

		println("\n- Installing", p.ID)
		parts := winget.BuildWingetInstallCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	println("\nDone!")
	helpers.PressAnyKeyOrWaitToExit()
}
