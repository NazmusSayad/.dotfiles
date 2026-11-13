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

type worktree struct {
	path   string
	branch string
}

func cleanBranches(yes bool) {
	currentBranch := helpers.GetCurrentGitBranchOrExit()

	topLevelOut, _ := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	currentPath := strings.TrimSpace(string(topLevelOut))

	worktreesOut, _ := exec.Command("git", "worktree", "list", "--porcelain").Output()
	var allWorktrees []worktree
	for _, block := range strings.Split(strings.TrimSpace(string(worktreesOut)), "\n\n") {
		var w worktree
		for _, line := range strings.Split(block, "\n") {
			line = strings.TrimSpace(line)
			if path, ok := strings.CutPrefix(line, "worktree "); ok {
				w.path = path
			}
			if branch, ok := strings.CutPrefix(line, "branch refs/heads/"); ok {
				w.branch = branch
			}
		}
		if w.path != "" {
			allWorktrees = append(allWorktrees, w)
		}
	}

	keptBranches := map[string]bool{}
	var worktrees []worktree
	for i, w := range allWorktrees {
		if i == 0 || w.path == currentPath {
			keptBranches[w.branch] = true
			continue
		}
		worktrees = append(worktrees, w)
	}

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
		if keptBranches[b] {
			continue
		}

		branches = append(branches, b)
	}

	if len(branches) == 0 && len(worktrees) == 0 {
		fmt.Println(aurora.Green("No other branches to delete"))
		return
	}

	worktreeBranches := map[string]bool{}
	for _, w := range worktrees {
		worktreeBranches[w.branch] = true
	}

	if len(branches) > 0 {
		colorfulBranches := []string{}
		for _, b := range branches {
			colorfulBranches = append(colorfulBranches, aurora.Red(string(b)).Bold().String())
		}

		fmt.Println(aurora.Yellow("Branches to delete:"), strings.Join(colorfulBranches, ", "))

		if !yes && !confirm("Delete these branches? ") {
			return
		}

		helpers.ExecNativeCommand([]string{"git", "worktree", "prune"})
		helpers.ExecNativeCommand([]string{"git", "prune", "--progress"})

		var freeBranches []string
		for _, b := range branches {
			if !worktreeBranches[b] {
				freeBranches = append(freeBranches, b)
			}
		}

		if len(freeBranches) > 0 {
			helpers.ExecNativeCommand(
				append([]string{"git", "branch", "-D"}, freeBranches...),
				helpers.ExecCommandOptions{
					Exit: true,
				},
			)
		}
	}

	if len(worktrees) == 0 {
		return
	}

	if len(branches) > 0 {
		fmt.Println()
	}

	colorfulWorktrees := []string{}
	for _, w := range worktrees {
		colorfulWorktrees = append(colorfulWorktrees, aurora.Red(w.path).Bold().String())
	}

	fmt.Println(aurora.Yellow("Worktrees to delete:"), strings.Join(colorfulWorktrees, ", "))

	if !yes && !confirm("Delete these worktrees? ") {
		return
	}

	for _, w := range worktrees {
		helpers.ExecNativeCommand(
			[]string{"git", "worktree", "remove", "--force", w.path},
			helpers.ExecCommandOptions{
				Exit: true,
			},
		)
	}

	var heldBranches []string
	for _, b := range branches {
		if worktreeBranches[b] {
			heldBranches = append(heldBranches, b)
		}
	}

	if len(heldBranches) > 0 {
		helpers.ExecNativeCommand(
			append([]string{"git", "branch", "-D"}, heldBranches...),
			helpers.ExecCommandOptions{
				Exit: true,
			},
		)
	}
}

func confirm(title string) bool {
	confirmed := true
	err := huh.NewConfirm().
		Title(title).
		Inline(true).
		Value(&confirmed).
		WithTheme(helpers.HuhTheme()).
		Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
		os.Exit(1)
	}

	return confirmed
}
