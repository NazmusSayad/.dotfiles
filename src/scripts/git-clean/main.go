package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"dotfiles/src/helpers"

	"charm.land/huh/v2"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	yes := false
	command := &cobra.Command{
		Use:   "git-clean",
		Short: "Delete all local branches except the current branch",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			cleanBranches(yes)
		},
	}
	command.Flags().BoolVarP(&yes, "yes", "y", false, "Delete branches without confirmation")

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func cleanBranches(yes bool) {
	currentBranch := helpers.GetCurrentGitBranchOrExit()

	branchesOut, _ := exec.Command("git", "branch", `--format=%(refname:short)`).Output()
	lines := strings.Split(strings.TrimRight(string(branchesOut), "\r\n"), "\n")

	var branches []string
	for _, b := range lines {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		if currentBranch != "" && strings.Contains(b, currentBranch) {
			continue
		}

		branches = append(branches, b)
	}

	if len(branches) == 0 {
		fmt.Println(aurora.Green("No other branches to delete"))
		return
	}

	colorfulBranches := []string{}
	for _, b := range branches {
		colorfulBranches = append(colorfulBranches, aurora.Red(string(b)).Bold().String())
	}

	fmt.Println(aurora.Yellow("Branches to delete:"), strings.Join(colorfulBranches, ", "))

	if !yes {
		confirmed := true
		err := huh.NewConfirm().
			Title("Delete these branches? ").
			Inline(true).
			Value(&confirmed).
			WithTheme(helpers.HuhTheme()).
			Run()
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return
			}
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
			os.Exit(1)
		}

		if !confirmed {
			return
		}
	}

	helpers.ExecNativeCommand([]string{"git", "prune", "--progress"})
	helpers.ExecNativeCommand(
		append([]string{"git", "branch", "-D"}, branches...),
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}
