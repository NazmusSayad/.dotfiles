package main

import (
	"fmt"
	"os"
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
		fmt.Fprintln(os.Stderr, "Usage: gpc [branch]")
		os.Exit(1)
	}

	ghUser := helpers.GetGitHubUser()
	fmt.Println(aurora.Faint("  GitHub user:"), ghUser)
	ghUser = ""

	remote := helpers.GetCurrentGitRemoteOrExit()
	remoteUrl := helpers.GetGitRemoteUrlOrExit(remote)

	branchCompare := ""
	if baseBranch != "" {
		branchCompare = baseBranch + "..." + targetBranch
		fmt.Println(
			aurora.Faint("  Merging"),
			aurora.Bold(baseBranch),
			aurora.Faint("into"),
			aurora.Bold(targetBranch),
		)
	} else {
		branchCompare = targetBranch
		fmt.Println(
			aurora.Faint("  Merging"),
			aurora.Bold(targetBranch),
			aurora.Faint("into"),
			aurora.Bold("default branch"),
		)
	}

	assignees := ""
	if ghUser != "" {
		assignees = "&assignees=" + ghUser
	}

	url := strings.Join([]string{remoteUrl + "/compare/" + branchCompare + "?expand=1", assignees}, "")
	fmt.Println(aurora.Faint("  " + url))

	helpers.Open(url, helpers.ExecCommandOptions{Exit: true})
}
