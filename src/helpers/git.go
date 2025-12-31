package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func GetCurrentGitBranchSafe() string {
	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetCurrentGitRemoteSafe() string {
	cmd := exec.Command("git", "remote")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetGitRemoteURLSafe(remote string) string {
	cmd := exec.Command("git", "remote", "get-url", remote)
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetCurrentGitBranch() string {
	branch := GetCurrentGitBranchSafe()
	if branch == "" {
		fmt.Println(aurora.Red("No branch found"))
		os.Exit(1)
	}

	return branch
}

func GetCurrentGitRemote() string {
	remote := GetCurrentGitRemoteSafe()
	if remote == "" {
		fmt.Println(aurora.Red("No remote found"))
		os.Exit(1)
	}

	return remote
}

func GetGitRemoteURL(remote string) string {
	url := GetGitRemoteURLSafe(remote)
	if url == "" {
		fmt.Println(aurora.Red("No remote URL found"))
		os.Exit(1)
	}

	return url
}
