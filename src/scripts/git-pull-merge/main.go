package main

import (
	"fmt"
	"os"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "git-pull-merge [branch]",
		Short: "Pull changes using merge and push the current branch",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			currentBranch := helpers.GetCurrentGitBranchOrExit()

			targetBranch := ""
			if len(args) == 0 {
				fmt.Println(aurora.Faint("No branch specified, using current branch"))
				targetBranch = currentBranch
			} else {
				targetBranch = args[0]
			}

			fmt.Printf(
				"Pulling changes from %s into %s (merge)\n", aurora.Yellow(targetBranch), aurora.Red(currentBranch),
			)

			remote := helpers.GetCurrentGitRemoteOrExit()

			helpers.ExecNativeCommand([]string{"git", "prune", "--progress"})
			helpers.ExecNativeCommand(
				[]string{"git", "pull", remote, targetBranch, "--progress", "--rebase=false"},
				helpers.ExecCommandOptions{
					Exit: true,
				},
			)

			helpers.ExecNativeCommand(
				[]string{"git", "push", remote, currentBranch, "--progress"},
				helpers.ExecCommandOptions{
					Exit: true,
				},
			)
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
