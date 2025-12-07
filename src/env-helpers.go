package helpers

import (
	"fmt"
	"strings"
)

type Scope string

const (
	ScopeUser    Scope = "User"
	ScopeMachine Scope = "Machine"
)

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
