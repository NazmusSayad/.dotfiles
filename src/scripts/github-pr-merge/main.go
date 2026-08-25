package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"dotfiles/src/helpers"

	gh "github.com/cli/go-gh/v2"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(
		os.Args[1:],
		"Usage: ghm [base-branch] [head-branch]",
	)

	pullRequestNumber, err := findPullRequestNumber(baseBranch, targetBranch)
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to find pull request"))
		os.Exit(1)
	}

	if pullRequestNumber == 0 {
		confirmed, err := helpers.CreateGithubPullRequest(baseBranch, targetBranch)
		if !confirmed {
			return
		}
		if err != nil {
			os.Exit(1)
		}

		pullRequestNumber, err = findPullRequestNumber(baseBranch, targetBranch)
		if err != nil {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to find pull request after creation"))
			os.Exit(1)
		}
		if pullRequestNumber == 0 {
			fmt.Fprintln(os.Stderr, aurora.Red("Pull request not found after creation"))
			os.Exit(1)
		}
	}

	fmt.Print(
		" Merge PR: ",
		aurora.Red(baseBranch).Bold(),
		aurora.Faint("<-"),
		aurora.Yellow(targetBranch).Bold(),
		aurora.Faint("[Press Enter]: "),
	)

	confirmation, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimRight(confirmation, "\r\n") != "" {
		fmt.Println(aurora.Red("Pull request merge cancelled"))
		return
	}

	if err := gh.ExecInteractive(
		context.Background(),
		"pr", "merge", strconv.Itoa(pullRequestNumber), "--merge",
	); err != nil {
		os.Exit(1)
	}
}

func findPullRequestNumber(baseBranch string, targetBranch string) (int, error) {
	pullRequestsOutput, _, err := gh.Exec(
		"pr",
		"list",
		"--state",
		"open",
		"--base",
		baseBranch,
		"--head",
		targetBranch,
		"--limit",
		"100",
		"--json",
		"number,baseRefName,headRefName",
	)
	if err != nil {
		return 0, err
	}

	var pullRequests []struct {
		Number      int    `json:"number"`
		BaseRefName string `json:"baseRefName"`
		HeadRefName string `json:"headRefName"`
	}
	if err := json.Unmarshal(pullRequestsOutput.Bytes(), &pullRequests); err != nil {
		return 0, err
	}

	pullRequestNumber := 0
	for _, pullRequest := range pullRequests {
		if pullRequest.BaseRefName != baseBranch || pullRequest.HeadRefName != targetBranch {
			continue
		}
		if pullRequestNumber != 0 {
			return 0, fmt.Errorf("multiple open pull requests found for %s <- %s", baseBranch, targetBranch)
		}
		pullRequestNumber = pullRequest.Number
	}

	return pullRequestNumber, nil
}
