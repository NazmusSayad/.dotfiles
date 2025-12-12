package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println(aurora.Red("This will reset the entire repository to the latest remote branch."))
	fmt.Println("Write 'yes' and press [Enter] to confirm.")
	fmt.Print("> ")

	confirm, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	confirm = strings.TrimRight(confirm, "\r\n")
	if confirm != "yes" {
		fmt.Println(aurora.Green("Reset aborted"))
		return
	}

	exec.Command("git", "fetch", "--all").Run()

	remoteURL := ""
	if out, err := exec.Command("git", "remote", "get-url", "origin").Output(); err == nil {
		remoteURL = strings.TrimSpace(string(out))
	}
	currentBranch := ""
	if out, err := exec.Command("git", "branch", "--show-current").Output(); err == nil {
		currentBranch = strings.TrimSpace(string(out))
	}

	fmt.Printf("> Branch: %s\n", currentBranch)
	fmt.Printf("> Remote: %s\n", remoteURL)

	remoteBranchesOut, _ := exec.Command("git", "branch", "-r", `--format=%(refname:short)`).Output()
	remoteBranches := strings.Split(strings.TrimRight(string(remoteBranchesOut), "\r\n"), "\n")
	for _, rb := range remoteBranches {
		rb = strings.TrimSpace(rb)
		if rb == "" {
			continue
		}
		i := strings.IndexByte(rb, '/')
		if i <= 0 || i == len(rb)-1 {
			continue
		}
		rb = rb[i+1:]
		if rb == currentBranch {
			continue
		}

		fmt.Printf("> Deleting remote branch: %s\n", rb)
		cmd := exec.Command("git", "push", "origin", "--delete", rb)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	fmt.Println(aurora.Red("> Deleting git folder..."))
	os.RemoveAll(".git")

	exec.Command("git", "init", "--initial-branch="+currentBranch).Run()
	exec.Command("git", "remote", "add", "origin", remoteURL).Run()
	exec.Command("git", "add", ".").Run()
	exec.Command("git", "commit", "-m", "Initial commit").Run()

	cmd := exec.Command("git", "push", "--force", "--set-upstream", "origin", currentBranch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}
