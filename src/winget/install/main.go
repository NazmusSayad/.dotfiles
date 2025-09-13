package main

import (
	"bufio"
	"dotfiles/src/winget"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	winget.ConfirmIsAdminExec()
	packages := winget.GetWingetPackages("./config/winget-apps.jsonc")
	fmt.Println("Installing packages, total:", len(packages))

	for _, p := range packages {
		fmt.Println("\n- Installing", p.ID)
		parts := winget.BuildWingetInstallCommands(p)
		cmd := exec.Command(parts[0], parts[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	}

	fmt.Println("\nDone!")
	fmt.Println("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	os.Exit(0)
}
