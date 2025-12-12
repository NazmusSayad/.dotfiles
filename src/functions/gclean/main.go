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
	current := ""
	if out, err := exec.Command("git", "branch", "--show-current").Output(); err == nil {
		current = strings.TrimSpace(string(out))
	}

	branchesOut, _ := exec.Command("git", "branch", `--format=%(refname:short)`).Output()
	lines := strings.Split(strings.TrimRight(string(branchesOut), "\r\n"), "\n")

	var branches []string
	for _, b := range lines {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		if current != "" && strings.Contains(b, current) {
			continue
		}

		branches = append(branches, b)
	}

	if len(branches) == 0 {
		fmt.Println(aurora.Green("No other branches to delete"))
		return
	}

	colorfulBranches := []string{}
	for _, b := range branches {
		colorfulBranches = append(colorfulBranches, aurora.Red(string(b)).Bold().String())
	}

	fmt.Println(aurora.Yellow("Branches to delete: "), strings.Join(colorfulBranches, ", "))
	fmt.Print(aurora.Faint("Press [Enter] to confirm, or any other key to cancel: "))

	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimRight(line, "\r\n") != "" {
		fmt.Println(aurora.Green("Cancelled branch deletion"))
		return
	}

	exec.Command("git", "prune", "--progress").Run()
	cmd := exec.Command("git", append([]string{"branch", "-D"}, branches...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			os.Exit(1)
		}
	}

	fmt.Println(aurora.Green("Branches deleted"))
}
