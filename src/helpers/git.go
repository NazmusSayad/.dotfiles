package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()

	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

func GetCurrentGitBranch() string {
	if !IsGitRepo() {
		return ""
	}

	cmd := exec.Command("git", "branch", "--show-current")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetCurrentGitRemote() string {
	if !IsGitRepo() {
		return ""
	}

	cmd := exec.Command("git", "remote")
	out, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetGitRemoteHTTPUrl(remote string) string {
	if !IsGitRepo() {
		return ""
	}

	cmd := exec.Command("git", "remote", "get-url", remote)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	url := strings.TrimSpace(string(out))

	if !strings.HasPrefix(url, "https://") {
		return ""
	}

	if result, ok := strings.CutSuffix(url, ".git"); ok {
		return result
	}

	return url
}

func InGitRepoOrExit() {
	if !IsGitRepo() {
		fmt.Println(aurora.Red("Not a git repository"))
		os.Exit(1)
	}
}

func GetCurrentGitBranchOrExit() string {
	InGitRepoOrExit()

	branch := GetCurrentGitBranch()
	if branch == "" {
		fmt.Println(aurora.Red("No branch found"))
		os.Exit(1)
	}

	return branch
}

func GetCurrentGitRemoteOrExit() string {
	InGitRepoOrExit()

	remote := GetCurrentGitRemote()
	if remote == "" {
		fmt.Println(aurora.Red("No remote found"))
		os.Exit(1)
	}

	return remote
}

func GetGitRemoteUrlOrExit(remote string) string {
	InGitRepoOrExit()

	url := GetGitRemoteHTTPUrl(remote)
	if url == "" {
		fmt.Println(aurora.Red("No remote URL found"))
		os.Exit(1)
	}

	return url
}
