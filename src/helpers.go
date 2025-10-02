package helpers

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Scope string

const (
	ScopeUser    Scope = "User"
	ScopeMachine Scope = "Machine"
)

func execPsCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("powershell command failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func ReadEnv(scope Scope, name string) (string, error) {
	return execPsCommand(
		fmt.Sprintf(`[System.Environment]::GetEnvironmentVariable("%s", [System.EnvironmentVariableTarget]::%s)`, name, scope),
	)
}

func WriteEnv(scope Scope, name, value string) (string, error) {
	return execPsCommand(
		fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("%s", "%s", [System.EnvironmentVariableTarget]::%s)`, name, value, scope),
	)
}

func AddToEnvPath(scope Scope, paths ...string) (string, error) {
	existingPath, err := ReadEnv(scope, "PATH")
	if err != nil {
		return "", err
	}

	existingPathArray := strings.Split(existingPath, ";")
	var filteredPaths []string
	for _, p := range existingPathArray {
		if p != "" {
			filteredPaths = append(filteredPaths, p)
		}
	}

	pathSet := make(map[string]bool)
	var uniquePaths []string

	for _, p := range filteredPaths {
		if !pathSet[p] {
			pathSet[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	for _, p := range paths {
		if !pathSet[p] {
			pathSet[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	newPath := strings.Join(uniquePaths, ";")
	return WriteEnv(scope, "PATH", newPath)
}

func EnsureAdminExecution() {
	psCmd := `if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	if err := cmd.Run(); err == nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "This program requires administrator privileges.")
		fmt.Fprintln(os.Stderr, "Trying to relaunch with elevated privileges...")

		exePath, e := os.Executable()
		if e != nil {
			_, _ = reader.ReadString('\n')
			os.Exit(1)
		}

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+exePath+"' -Verb RunAs")
		if err := cmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to relaunch with elevated privileges. Press Enter to exit...")
			_, _ = reader.ReadString('\n')
			os.Exit(1)
		}

		os.Exit(0)
	}
}

func WaitForInputAndExit() {
	fmt.Println("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
	os.Exit(0)
}
