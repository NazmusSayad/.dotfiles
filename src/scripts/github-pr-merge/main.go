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
		"1",
		"--json",
		"number,mergeable",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to find pull request"))
		os.Exit(1)
	}

	var pullRequests []struct {
		Number    int    `json:"number"`
		Mergeable string `json:"mergeable"`
	}
	if err := json.Unmarshal(pullRequestsOutput.Bytes(), &pullRequests); err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to read pull request"))
		os.Exit(1)
	}
	if len(pullRequests) == 0 {
		fmt.Fprintln(os.Stderr, aurora.Red("Open pull request not found"))
		os.Exit(1)
	}

	pullRequest := pullRequests[0]
	if pullRequest.Mergeable == "CONFLICTING" {
		fmt.Fprintln(os.Stderr, aurora.Red("Pull request has merge conflicts"))
		os.Exit(1)
	}
	if pullRequest.Mergeable != "MERGEABLE" {
		fmt.Fprintln(os.Stderr, aurora.Red("Pull request mergeability is unknown"))
		os.Exit(1)
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
		"pr", "merge", strconv.Itoa(pullRequest.Number), "--merge",
	); err != nil {
		os.Exit(1)
	}
}
