package helpers

import (
	"os/exec"
	"strings"
)

func GetCurrentGitRemote() string {
	cmd := exec.Command("git", "remote")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func getGitRemoteURL(remote string) string {
	cmd := exec.Command("git", "remote", "get-url", remote)
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
