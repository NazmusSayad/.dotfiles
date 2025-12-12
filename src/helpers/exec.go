package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func sudoAvailable() bool {
	_, err := exec.LookPath("sudo")
	return err == nil
}

func isRunningAsAdmin() bool {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"(New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to check if running as admin.")
		os.Exit(1)
	}

	return strings.TrimSpace(out.String()) == "True"
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

	if sudoAvailable() {
		cmd := exec.Command("sudo", exe)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Failed to run sudo.")
			os.Exit(1)
		}

		os.Exit(0)
	}

	fmt.Println("Relaunching with elevated privileges...")

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+exe+"' -Verb RunAs")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to relaunch with elevated privileges.")
		os.Exit(1)
	}

	os.Exit(0)
}
