package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if !helpers.IsGitRepo() {
		fmt.Fprintln(os.Stderr, aurora.Red("Not a git repository"))
		os.Exit(1)
	}

	baseBranch := ""
	targetBranch := helpers.GetCurrentGitBranchOrExit()

	if len(os.Args) == 1 {
	} else if len(os.Args) == 2 {
		baseBranch = os.Args[1]
	} else if len(os.Args) == 3 {
		baseBranch = os.Args[1]
		targetBranch = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "Usage: gpc [base-branch] [head-branch]")
		os.Exit(1)
	}

	if baseBranch == "" {
		defaultBranchOutput, err := exec.Command(
			"gh",
			"repo",
			"view",
			"--json",
			"defaultBranchRef",
			"--jq",
			".defaultBranchRef.name",
		).Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to resolve default branch"))
			os.Exit(1)
		}

		baseBranch = strings.TrimSpace(string(defaultBranchOutput))
		if baseBranch == "" {
			fmt.Fprintln(os.Stderr, aurora.Red("Default branch not found"))
			os.Exit(1)
		}
	}

	fmt.Print(
		" Create PR: ",
		aurora.Red(baseBranch).Bold(),
		aurora.Faint("<-"),
		aurora.Yellow(targetBranch).Bold(),
		aurora.Faint("[Press Enter]: "),
	)

	confirmation, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimRight(confirmation, "\r\n") != "" {
		fmt.Println(aurora.Red("Pull request creation cancelled"))
		return
	}

	helpers.ExecNativeCommand(
		[]string{
			"gh", "pr", "create", "--fill",
			"--assignee", "@me",
			"--base", baseBranch,
			"--head", targetBranch,
		},
		helpers.ExecCommandOptions{Exit: true},
	)
}
