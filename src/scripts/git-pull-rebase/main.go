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
		Use:   "git-pull-rebase [branch]",
		Short: "Pull changes using rebase and push the current branch",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			pullRebase(args)
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func pullRebase(args []string) {
	currentBranch := helpers.GetCurrentGitBranchOrExit()

	targetBranch := ""
	if len(args) == 0 {
		fmt.Println(aurora.Faint("No branch specified, using current branch"))
		targetBranch = currentBranch
	} else {
		targetBranch = args[0]
	}

	fmt.Printf(
		"Pulling changes from %s into %s (rebase)\n", aurora.Yellow(targetBranch), aurora.Red(currentBranch),
	)

	remote := helpers.GetCurrentGitRemoteOrExit()

	helpers.ExecNativeCommand([]string{"git", "prune", "--progress"})
	helpers.ExecNativeCommand(
		[]string{"git", "pull", remote, targetBranch, "--progress", "--rebase=true"},
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
}
