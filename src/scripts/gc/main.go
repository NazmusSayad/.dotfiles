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
		helpers.ExecNativeCommand(helpers.ExecCommandOptions{
			Command: "git",
			Args:    append([]string{"checkout"}, os.Args[1:]...),
			Exit:    true,
		})
	} else {
		helpers.ExecNativeCommand(helpers.ExecCommandOptions{
			Command: "git",
			Args:    append([]string{"checkout", "-b"}, os.Args[1:]...),
			Exit:    true,
		})
	}
}

func isLocalBranchExists(branch string) bool {
	return helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "git",
		Args:    []string{"rev-parse", "--verify", "--quiet", "refs/heads/" + branch},
	}) == nil
}

func isRemoteBranchExists(remote string, branch string) bool {
	return helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "git",
		Args:    []string{"rev-parse", "--verify", "--quiet", "refs/remotes/" + remote + "/" + branch},
	}) == nil
}
