package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"dotfiles/src/helpers"

	"charm.land/huh/v2"
	gh "github.com/cli/go-gh/v2"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "github-pr-merge [base-branch] [head-branch]",
		Short: "Merge a GitHub pull request",
		Args:  cobra.MaximumNArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(args)

			pullRequest, err := helpers.FindGithubPullRequest(baseBranch, targetBranch)
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

				pullRequest, err = helpers.FindGithubPullRequest(baseBranch, targetBranch)
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

			confirmed := true
			err = huh.NewConfirm().
				Title(fmt.Sprint(
					aurora.Magenta(" Merge PR "),
					aurora.Cyan("#"+strconv.Itoa(pullRequest.Number)).Bold(), ": ",
					aurora.Red(baseBranch),
					aurora.Faint("<-"),
					aurora.Yellow(targetBranch),
				)).
				Inline(true).
				Value(&confirmed).
				WithTheme(huh.ThemeFunc(huh.ThemeBase)).
				Run()
			if err != nil {
				if errors.Is(err, huh.ErrUserAborted) {
					fmt.Println(aurora.Red("Pull request merge cancelled"))
					return
				}
				fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
				os.Exit(1)
			}

			if !confirmed {
				fmt.Println(aurora.Red("Pull request merge cancelled"))
				return
			}

			if err := gh.ExecInteractive(
				context.Background(),
				"pr", "merge", strconv.Itoa(pullRequest.Number), "--merge",
			); err != nil {
				os.Exit(1)
			}
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
