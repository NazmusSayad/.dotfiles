package main

import (
	"fmt"
	"os"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	web := false
	yes := false
	command := &cobra.Command{
		Use:   "github-pr-create [base-branch] [head-branch]",
		Short: "Create a GitHub pull request",
		Args:  cobra.MaximumNArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(args)

			created, err := helpers.CreateGithubPullRequest(baseBranch, targetBranch, yes)
			if err != nil && !web {
				os.Exit(1)
			}
			if web && (created || err != nil) {
				pullRequest, err := helpers.FindGithubPullRequest(baseBranch, targetBranch)
				if err != nil {
					fmt.Fprintln(os.Stderr, aurora.Red("Failed to find pull request"))
					os.Exit(1)
				}
				if pullRequest.Number == 0 {
					fmt.Fprintln(os.Stderr, aurora.Red("Pull request not found"))
					os.Exit(1)
				}

				if err := helpers.Open(pullRequest.URL); err != nil {
					os.Exit(1)
				}
			}
		},
	}
	command.Flags().BoolVarP(&web, "web", "w", false, "Open the pull request in a browser")
	command.Flags().BoolVarP(&yes, "yes", "y", false, "Create the pull request without confirmation")

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
