package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(cmd string, args ...string) (string, int) {
	var b bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Stdout = &b
	c.Stderr = &b
	err := c.Run()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return b.String(), e.ExitCode()
		}
		return b.String(), 1
	}
	return b.String(), 0
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "❌ Branch name required")
		os.Exit(1)
	}
	branch := args[0]
	if strings.HasPrefix(branch, "-") {
		fmt.Fprintln(os.Stderr, "❌ Invalid branch name:", branch)
		os.Exit(1)
	}
	remoteOut, _ := run("git", "remote")
	remote := ""
	for _, r := range strings.Split(strings.TrimSpace(remoteOut), "\n") {
		if strings.TrimSpace(r) != "" {
			remote = strings.Fields(r)[0]
			break
		}
	}
	existsLocal := false
	if _, code := run("git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch); code == 0 {
		existsLocal = true
	}
	existsRemote := false
	if remote != "" {
		if _, code := run("git", "rev-parse", "--verify", "--quiet", "refs/remotes/"+remote+"/"+branch); code == 0 {
			existsRemote = true
		}
	}
	if existsLocal || existsRemote {
		checkoutArgs := append([]string{"checkout"}, args...)
		exec.Command("git", checkoutArgs...).Run()
		return
	}
	newArgs := append([]string{"checkout", "-b"}, args...)
	exec.Command("git", newArgs...).Run()
}

