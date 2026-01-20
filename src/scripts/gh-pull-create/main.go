package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
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

	targetBranch := ""
	if len(os.Args) == 1 {
		currentBranch := helpers.GetCurrentGitBranchOrExit()
		targetBranch = currentBranch

		fmt.Println("Using current branch:", aurora.Yellow(currentBranch).String())
	} else if len(os.Args) == 2 {
		targetBranch = os.Args[1]
	} else {
		fmt.Fprintln(os.Stderr, "Usage: ghp [branch]")
		os.Exit(1)
	}

	ghUser := helpers.GetGitHubUser()
	remote := helpers.GetCurrentGitRemoteOrExit()
	remoteUrl := helpers.GetGitRemoteUrlOrExit(remote)

	url := strings.Join([]string{remoteUrl + "/compare/" + targetBranch + "?expand=1", utils.Ternary(ghUser != "", "&assignees="+ghUser, "")}, "")
	fmt.Println(aurora.Faint("URL: " + url))

	helpers.ExecNativeCommand(
		[]string{"rundll32", "url.dll,FileProtocolHandler", url},
		helpers.ExecCommandOptions{Exit: true},
	)
}
