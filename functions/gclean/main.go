package main

import (
	"bufio"
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
	listOut, _ := run("git", "branch", "--format=%(refname:short)")
	var branches []string
	for _, b := range strings.Split(listOut, "\n") {
		name := strings.TrimSpace(b)
		if name != "" && name != current {
			branches = append(branches, name)
		}
	}
	if len(branches) == 0 {
		fmt.Println("No other branches to delete")
		return
	}
	fmt.Print("Branches to delete: ", strings.Join(branches, ", "), "\n")
	fmt.Print("Press [Enter] to confirm, or any other key to cancel: ")
	r := bufio.NewReader(os.Stdin)
	input, _ := r.ReadString('\n')
	if input == "\n" {
		run("git", "prune", "--progress")
		args := append([]string{"branch", "-D"}, branches...)
		run("git", args...)
		return
	}
	fmt.Println("Cancelled branch deletion")
}

