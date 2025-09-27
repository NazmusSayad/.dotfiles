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

	updatablePackages := winget.GetUpdatablePackageIDs()
	fmt.Println("Found", len(updatablePackages), "packages with available updates")

	for _, p := range packages {
		if p.IgnoreUpgrade || p.Version != "" {
			fmt.Println("\n- Skipping", p.ID)
			continue
		}

		if !slices.Contains(updatablePackages, p.ID) {
			fmt.Println("\n- No updates available for", p.ID)
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
	helpers.WaitForInputAndExit()
}
