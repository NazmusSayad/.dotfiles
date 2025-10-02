package main

import (
	helpers "dotfiles/src"
	"dotfiles/src/winget"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	packages := winget.GetWingetPackages("./config/winget-apps.jsonc")
	fmt.Println("Installing packages, total:", len(packages))

	for _, p := range packages {
		if p.IgnoreInstall {
			fmt.Println("\n- Skipping", p.ID)
			continue
		}

		fmt.Println("\n- Installing", p.ID)
		parts := winget.BuildWingetInstallCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	fmt.Println("\nDone!")
	helpers.WaitForInputAndExit()
}
