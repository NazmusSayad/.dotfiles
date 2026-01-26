package helpers

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

type Scope string

const (
	ScopeUser    Scope = "User"
	ScopeMachine Scope = "Machine"
)

func execPsCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("powershell command failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func psEscape(s string) string {
	return strings.ReplaceAll(s, `'`, `''`)
}

func ReadEnv(scope Scope, name string) (string, error) {
	return execPsCommand(
		fmt.Sprintf(
			"[System.Environment]::GetEnvironmentVariable('%s', [System.EnvironmentVariableTarget]::%s)",
			psEscape(name),
			scope,
		),
	)
}

func WriteEnv(scope Scope, name, value string) (string, error) {
	return execPsCommand(
		fmt.Sprintf(
			"[System.Environment]::SetEnvironmentVariable('%s', '%s', [System.EnvironmentVariableTarget]::%s)",
			psEscape(name),
			psEscape(value),
			scope,
		),
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
