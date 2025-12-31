package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	currentBranch := ""
	if out, err := exec.Command("git", "branch", "--show-current").Output(); err == nil {
		currentBranch = strings.TrimSpace(string(out))
	}

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
		"Pulling changes from %s into %s (default)\n", aurora.Yellow(targetBranch), aurora.Red(currentBranch),
	)

	remote := helpers.GetCurrentGitRemote()
	if remote == "" {
		fmt.Println(aurora.Red("No remote found"))
		os.Exit(1)
	}

	helpers.ExecNativeCommand([]string{"git", "prune", "--progress"})
	helpers.ExecNativeCommand(
		[]string{"git", "pull", remote, targetBranch, "--progress", "--rebase=false"},
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}
