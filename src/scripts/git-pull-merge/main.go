package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	currentBranch := helpers.GetCurrentGitBranchOrExit()

	targetBranch := ""
	if len(os.Args) == 1 {
		fmt.Println(aurora.Faint("No branch specified, using current branch"))
		targetBranch = currentBranch
	} else if len(os.Args) == 2 {
		targetBranch = os.Args[1]
	} else {
		fmt.Fprintln(os.Stderr, "Usage: gp [branch]")
		os.Exit(1)
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
}
