package main

import (
	"dotfiles/src/helpers/winget"
	"os"
	"os/exec"
	"slices"
)

func main() {
	packages := winget.GetWingetPackages()

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

		if p.SkipUpgrade || p.Version != "" {
			println("\n- Skipping", p.ID)
			continue
		}

		println("\n- Upgrading", p.ID)
		parts := winget.BuildWingetUpgradeCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	println("\nDone!")
}
