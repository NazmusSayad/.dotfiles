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
	currentOut, _ := run("git", "branch", "--show-current")
	current := strings.TrimSpace(currentOut)
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "Usage: gp [branch]")
		os.Exit(1)
	}
	target := current
	if len(args) == 1 {
		target = args[0]
	} else {
		fmt.Println("No branch specified, using current branch")
	}
	fmt.Print("Pulling changes from ", target, " into ", current, " (rebase)\n")
	run("git", "prune", "--progress")
	run("git", "pull", "origin", target, "--progress", "--rebase")
}
