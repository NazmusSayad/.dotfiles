package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintln(os.Stderr, "Usage: ghp [branch]")
		os.Exit(1)
	}

	ghUser := helpers.GetGitHubUser()
	remote := helpers.GetCurrentGitRemoteOrExit()
	remoteUrl := helpers.GetGitRemoteUrlOrExit(remote)

	branchCompare := ""
	if baseBranch != "" {
		branchCompare = baseBranch + "..." + targetBranch
	} else {
		branchCompare = targetBranch
	}

	assignees := ""
	if ghUser != "" {
		assignees = "&assignees=" + ghUser
	}

	url := strings.Join([]string{remoteUrl + "/compare/" + branchCompare + "?expand=1", assignees}, "")
	fmt.Println(aurora.Faint("  " + url))

	helpers.ExecNativeCommand(
		[]string{"rundll32", "url.dll,FileProtocolHandler", url},
		helpers.ExecCommandOptions{Exit: true},
	)
}
