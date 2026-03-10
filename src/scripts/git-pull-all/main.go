package main

import (
	"fmt"
	"os/exec"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println(aurora.Yellow("Pulling changes from all branches:"), strings.Join(getGitBranches(), ", "))

	helpers.ExecNativeCommand(
		[]string{"git", "pull", "--all"},
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
