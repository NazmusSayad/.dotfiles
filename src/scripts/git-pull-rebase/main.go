package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	currentBranch := helpers.GetCurrentGitBranch()

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
		"Pulling changes from %s into %s (rebase)\n", aurora.Yellow(targetBranch), aurora.Red(currentBranch),
	)

	remote := helpers.GetCurrentGitRemote()

	helpers.ExecNativeCommand([]string{"git", "prune", "--progress"})
	helpers.ExecNativeCommand(
		[]string{"git", "pull", remote, targetBranch, "--progress", "--rebase=true"},
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}
