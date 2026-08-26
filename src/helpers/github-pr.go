package helpers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	gh "github.com/cli/go-gh/v2"
	"github.com/logrusorgru/aurora/v4"
)

func GetGithubPullRequestBranchesOrExit(args []string, usage string) (string, string) {
	baseBranch := ""
	targetBranch := GetCurrentGitBranchOrExit()

	if len(args) == 1 {
		baseBranch = args[0]
	} else if len(args) == 2 {
		baseBranch = args[0]
		targetBranch = args[1]
	} else if len(args) > 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	if baseBranch != "" {
		return baseBranch, targetBranch
	}

	defaultBranchOutput, _, err := gh.Exec(
		"repo",
		"view",
		"--json",
		"defaultBranchRef",
		"--jq",
		".defaultBranchRef.name",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to resolve default branch"))
		os.Exit(1)
	}

	baseBranch = strings.TrimSpace(defaultBranchOutput.String())
	if baseBranch == "" {
		fmt.Fprintln(os.Stderr, aurora.Red("Default branch not found"))
		os.Exit(1)
	}

	return baseBranch, targetBranch
}

func CreateGithubPullRequest(baseBranch string, targetBranch string) (bool, error) {
	fmt.Print(
		" Create PR: ",
		aurora.Red(baseBranch).Bold(),
		aurora.Faint("<-"),
		aurora.Yellow(targetBranch).Bold(),
		aurora.Faint("[Press Enter]: "),
	)

	confirmation, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimRight(confirmation, "\r\n") != "" {
		fmt.Println(aurora.Red("Pull request creation cancelled"))
		return false, nil
	}

	err := gh.ExecInteractive(
		context.Background(),
		"pr", "create", "--fill",
		"--assignee", "@me",
		"--base", baseBranch,
		"--head", targetBranch,
	)
	return true, err
}
