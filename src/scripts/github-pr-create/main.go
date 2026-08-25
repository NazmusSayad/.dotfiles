package main

import (
	"os"

	"dotfiles/src/helpers"
)

func main() {
	baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(
		os.Args[1:],
		"Usage: ghp [base-branch] [head-branch]",
	)

	_, err := helpers.CreateGithubPullRequest(baseBranch, targetBranch)
	if err != nil {
		os.Exit(1)
	}
}
