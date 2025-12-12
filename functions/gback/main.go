package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
		fmt.Fprintln(os.Stderr, "❌ Commit hash required")
		os.Exit(1)
	}
	run("git", "restore", "--source", args[0], "--", ".")
}
