package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"dotfiles/src/helpers"

	gh "github.com/cli/go-gh/v2"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	baseBranch, targetBranch := helpers.GetGithubPullRequestBranchesOrExit(
		os.Args[1:],
		"Usage: ghp [base-branch] [head-branch]",
	)

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

	if err := gh.ExecInteractive(
		context.Background(),
		"pr", "create", "--fill",
		"--assignee", "@me",
		"--base", baseBranch,
		"--head", "refs/heads/"+targetBranch,
	); err != nil {
		os.Exit(1)
	}
}
