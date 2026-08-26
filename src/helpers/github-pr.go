package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"charm.land/huh/v2"
	gh "github.com/cli/go-gh/v2"
	"github.com/logrusorgru/aurora/v4"
)

func GetGithubPullRequestBranchesOrExit(args []string) (string, string) {
	baseBranch := ""
	targetBranch := GetCurrentGitBranchOrExit()

	if len(args) == 1 {
		baseBranch = args[0]
	} else if len(args) == 2 {
		baseBranch = args[0]
		targetBranch = args[1]
	}

	if baseBranch != "" {
		return baseBranch, targetBranch
	}

	defaultBranchOutput, _, err := gh.Exec(
		"repo",
		"view",
		"--json",
		"defaultBranchRef",
		"--jq",
		".defaultBranchRef.name",
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to resolve default branch"))
		os.Exit(1)
	}

	baseBranch = strings.TrimSpace(defaultBranchOutput.String())
	if baseBranch == "" {
		fmt.Fprintln(os.Stderr, aurora.Red("Default branch not found"))
		os.Exit(1)
	}

	return baseBranch, targetBranch
}

func CreateGithubPullRequest(baseBranch string, targetBranch string) (bool, error) {
	title := fmt.Sprint(
		aurora.Green(" Create PR").String()+": ",
		aurora.Red(baseBranch).Bold(),
		aurora.Faint("<-"),
		aurora.Yellow(targetBranch).Bold(),
	)
	confirmed := true
	err := huh.NewConfirm().
		Title(title + " ").
		Inline(true).
		Value(&confirmed).
		WithTheme(huh.ThemeFunc(huh.ThemeBase)).
		Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(aurora.Red("Pull request creation cancelled"))
			return false, nil
		}
		return false, err
	}
	if !confirmed {
		fmt.Println(aurora.Red("Pull request creation cancelled"))
		return false, nil
	}

	err = gh.ExecInteractive(
		context.Background(),
		"pr", "create", "--fill",
		"--assignee", "@me",
		"--base", baseBranch,
		"--head", targetBranch,
	)
	return true, err
}
