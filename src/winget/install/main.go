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

	fmt.Println("Installing packages, total:", len(config.Packages))

	for _, p := range config.Packages {
		fmt.Println("\n- Installing", p.ID)
		parts := winget.BuildWingetInstallCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	fmt.Println("\nDone!")
}
