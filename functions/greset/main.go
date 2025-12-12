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
	fmt.Println("This will reset the entire repository to the latest remote branch.")
	fmt.Println("Write 'yes' and press [Enter] to confirm.")
	r := bufio.NewReader(os.Stdin)
	confirm, _ := r.ReadString('\n')
	if strings.TrimSpace(confirm) != "yes" {
		fmt.Println("Reset aborted")
		return
	}
	run("git", "fetch", "--all")
	remoteOut, _ := run("git", "remote", "get-url", "origin")
	remoteURL := strings.TrimSpace(remoteOut)
	currentOut, _ := run("git", "branch", "--show-current")
	current := strings.TrimSpace(currentOut)
	fmt.Print("> Branch: ", current, "\n")
	fmt.Print("> Remote: ", remoteURL, "\n")
	rbOut, _ := run("git", "branch", "-r", "--format=%(refname:short)")
	for _, line := range strings.Split(rbOut, "\n") {
		name := strings.TrimSpace(line)
		if name == "" || !strings.Contains(name, "/") {
			continue
		}
		parts := strings.SplitN(name, "/", 2)
		if len(parts) != 2 {
			continue
		}
		branch := parts[1]
		if branch == current {
			continue
		}
		fmt.Print("> Deleting remote branch: ", branch, "\n")
		run("git", "push", "origin", "--delete", branch)
	}
	fmt.Println("> Deleting git folder...")
	os.RemoveAll(".git")
	exec.Command("git", "init", "--initial-branch="+current).Run()
	if remoteURL != "" {
		exec.Command("git", "remote", "add", "origin", remoteURL).Run()
	}
	run("git", "add", ".")
	run("git", "commit", "-m", "Initial commit")
	run("git", "push", "--force", "--set-upstream", "origin", current)
}

