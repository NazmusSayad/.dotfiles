package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func execPsCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("powershell command failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func sudoAvailable() bool {
	_, err := exec.LookPath("sudo")
	return err == nil
}

func IsRunningAsAdmin() bool {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"(New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		println("Failed to check if running as admin.")
		os.Exit(1)
	}

	return strings.TrimSpace(out.String()) == "True"
}

func EnsureAdminExecution() {
	if IsRunningAsAdmin() {
		return
	}

	exePath, exePathErr := os.Executable()
	if exePathErr != nil {
		println("Failed to get executable path.")
		os.Exit(1)
	}

	if sudoAvailable() {
		cmd := exec.Command("sudo", exePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			println("Failed to run sudo.")
			os.Exit(1)
		}

		os.Exit(1)
		return
	}

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+exePath+"' -Verb RunAs")
	if err := cmd.Run(); err != nil {
		println("Failed to relaunch with elevated privileges.")
		println("Press Enter to exit...")
		os.Exit(1)
	}

	println("Relaunched with elevated privileges.")
	os.Exit(0)
}
