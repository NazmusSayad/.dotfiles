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

type pullRequest struct {
	Number      int    `json:"number"`
	URL         string `json:"url"`
	BaseRefName string `json:"baseRefName"`
	HeadRefName string `json:"headRefName"`
}

func main() {
	baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(
		os.Args[1:],
		"Usage: ghm [base-branch] [head-branch]",
	)

	pullRequest, err := findPullRequest(baseBranch, targetBranch)
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to find pull request"))
		os.Exit(1)
	}

	createdPullRequest := false
	if pullRequest.Number == 0 {
		confirmed, err := helpers.CreateGithubPullRequest(baseBranch, targetBranch)
		if !confirmed {
			return
		}
		if err != nil {
			os.Exit(1)
		}

		pullRequest, err = findPullRequest(baseBranch, targetBranch)
		if err != nil {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to find pull request after creation"))
			os.Exit(1)
		}
		if pullRequest.Number == 0 {
			fmt.Fprintln(os.Stderr, aurora.Red("Pull request not found after creation"))
			os.Exit(1)
		}
		createdPullRequest = true
	}

	if !createdPullRequest {
		fmt.Println(pullRequest.URL)
	}

	fmt.Print(
		"󰽜 Merge PR #",
		aurora.Cyan(pullRequest.Number).Bold(),
		": ",
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
		"pr", "merge", strconv.Itoa(pullRequest.Number), "--merge",
	); err != nil {
		os.Exit(1)
	}
}

func findPullRequest(baseBranch string, targetBranch string) (pullRequest, error) {
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
		"number,url,baseRefName,headRefName",
	)
	if err != nil {
		return pullRequest{}, err
	}

	var pullRequests []pullRequest
	if err := json.Unmarshal(pullRequestsOutput.Bytes(), &pullRequests); err != nil {
		return pullRequest{}, err
	}

	var matchingPullRequest pullRequest
	for _, candidate := range pullRequests {
		if candidate.BaseRefName != baseBranch || candidate.HeadRefName != targetBranch {
			continue
		}
		if matchingPullRequest.Number != 0 {
			return pullRequest{}, fmt.Errorf("multiple open pull requests found for %s <- %s", baseBranch, targetBranch)
		}
		matchingPullRequest = candidate
	}

	return matchingPullRequest, nil
}
