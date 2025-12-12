package main

import (
	"fmt"
	"io"
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

	localRef := exec.Command("git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch)
	localRef.Stdout = io.Discard
	localRef.Stderr = io.Discard
	localOk := localRef.Run() == nil

	remoteRef := exec.Command("git", "rev-parse", "--verify", "--quiet", "refs/remotes/"+remote+"/"+branch)
	remoteRef.Stdout = io.Discard
	remoteRef.Stderr = io.Discard
	remoteOk := remoteRef.Run() == nil

	var cmd *exec.Cmd
	if localOk || remoteOk {
		cmd = exec.Command("git", append([]string{"checkout"}, os.Args[1:]...)...)
	} else {
		cmd = exec.Command("git", append([]string{"checkout", "-b"}, os.Args[1:]...)...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}
