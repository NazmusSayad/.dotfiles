//go:build !windows

package helpers

import (
	"fmt"
	"os"
	"os/exec"

	"dotfiles/src/utils"
)

func isRunningAsAdmin() bool {
	return os.Geteuid() == 0
}

func EnsureAdminExecution() {
	if isRunningAsAdmin() {
		return
	}

	exe, exeErr := os.Executable()
	if exeErr != nil {
		fmt.Println("Failed to get executable path.")
		os.Exit(1)
	}

	if utils.IsCommandInPath("sudo") {
		cmd := exec.Command("sudo", exe)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to run sudo as admin.")
			os.Exit(1)
		}

		os.Exit(0)
	}

	fmt.Println("sudo not found in PATH")
	os.Exit(1)
}
