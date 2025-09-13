package main

import (
	"dotfiles/config"
	"dotfiles/src/winget"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	winget.ConfirmIsAdminExec()

	if len(config.Packages) == 0 {
		fmt.Fprintln(os.Stderr, "no packages configured")
		os.Exit(1)
	}

	fmt.Println("Upgrading packages, total:", len(config.Packages))

	for _, p := range config.Packages {
		parts := winget.BuildWingetUpgradeCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			}

			fmt.Fprintln(os.Stderr, "failed to upgrade", p.ID, ":", err)
			os.Exit(1)
		}
	}
}
