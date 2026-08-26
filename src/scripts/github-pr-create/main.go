package main

import (
	"os"

	"dotfiles/src/helpers"

	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "github-pr-create [base-branch] [head-branch]",
		Short: "Create a GitHub pull request",
		Args:  cobra.MaximumNArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(args)

			_, err := helpers.CreateGithubPullRequest(baseBranch, targetBranch)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
