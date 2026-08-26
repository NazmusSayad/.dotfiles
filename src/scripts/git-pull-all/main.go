package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "git-pull-all",
		Short: "Pull and push changes for all branches",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			pullAll()
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

func pullAll() {
	fmt.Println(aurora.Yellow("Pulling changes from all branches:"), strings.Join(getGitBranches(), ", "))

	helpers.ExecNativeCommand(
		[]string{"git", "pull", "--all"},
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)

	helpers.ExecNativeCommand(
		[]string{"git", "push", "--progress"},
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}

func getGitBranches() []string {
	branches := []string{}
	branchesOut, _ := exec.Command("git", "branch", `--format=%(refname:short)`).Output()
	lines := strings.SplitSeq(strings.TrimRight(string(branchesOut), "\r\n"), "\n")

	for b := range lines {
		b = strings.TrimSpace(b)
		if b != "" {
			branches = append(branches, aurora.Red(b).Bold().String())
		}
	}

	return branches
}
