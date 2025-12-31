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
