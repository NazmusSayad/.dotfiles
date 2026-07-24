package helpers

import (
	"os/exec"
	"strings"
)

type GhHostInfo struct {
	Host   string `json:"host"`
	Login  string `json:"login"`
	Active bool   `json:"active"`
}

type GhAuthStatus struct {
	Hosts map[string][]GhHostInfo `json:"hosts"`
}

func GetGitHubToken() string {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}
