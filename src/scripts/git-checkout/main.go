package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	branch := ""
	if len(os.Args) > 1 {
		branch = os.Args[1]
	}

	if branch == "" {
		fmt.Println(aurora.Red("Branch name required"))
		os.Exit(1)
	}

	if strings.HasPrefix(branch, "-") {
		fmt.Println(aurora.Red("Invalid branch name: " + branch))
		os.Exit(1)
	}

	remote := ""
	if out, err := exec.Command("git", "remote").Output(); err == nil {
		s := strings.TrimSpace(string(out))
		if s != "" {
			remote = strings.Split(s, "\n")[0]
		}
	}

	if isLocalBranchExists(branch) || isRemoteBranchExists(remote, branch) {
		args := append([]string{"git", "checkout"}, os.Args[1:]...)
		helpers.ExecNativeCommand(args, helpers.ExecCommandOptions{Exit: true})
	} else {
		args := append([]string{"git", "checkout", "-b"}, os.Args[1:]...)
		helpers.ExecNativeCommand(args, helpers.ExecCommandOptions{Exit: true})
	}
}

func isLocalBranchExists(branch string) bool {
	return helpers.ExecNativeCommand(
		[]string{"git", "rev-parse", "--verify", "--quiet", "refs/heads/" + branch},
		helpers.ExecCommandOptions{
			Silent: true,
		},
	) == nil
}

func isRemoteBranchExists(remote string, branch string) bool {
	return helpers.ExecNativeCommand(
		[]string{"git", "rev-parse", "--verify", "--quiet", "refs/remotes/" + remote + "/" + branch},
		helpers.ExecCommandOptions{
			Silent: true,
		},
	) == nil
}
