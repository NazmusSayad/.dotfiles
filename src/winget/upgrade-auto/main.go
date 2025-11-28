package main

import (
	helpers "dotfiles/src"
	"dotfiles/src/winget"
	"fmt"
	"os"
	"os/exec"
	"slices"
)

func main() {
	packages := winget.GetWingetPackages("./config/winget-apps.jsonc")

	var packageIDs []string
	upgradeablePackages := winget.GetUpgradeablePackages()

	for _, p := range upgradeablePackages {
		println("")
		println("ID:", p.ID)
		println("Current Version:", p.Version)
		println("Available Version:", p.Available)
		packageIDs = append(packageIDs, p.ID)
	}

	for _, p := range packages {
		if !slices.Contains(packageIDs, p.ID) {
			continue
		}

		if p.IgnoreUpgrade || p.Version != "" {
			fmt.Println("\n- Skipping", p.ID)
			continue
		}

		fmt.Println("\n- Upgrading", p.ID)
		parts := winget.BuildWingetUpgradeCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	fmt.Println("\nDone!")
	helpers.PressAnyKeyOrWaitToExit()
}
