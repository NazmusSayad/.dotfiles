package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(aurora.Red("Usage: c <repository-path> [additional-arguments]"))
		os.Exit(1)
	}

	inputPath := os.Args[1]
	resolvedPath := ""

	re := regexp.MustCompile(`^[^/]+(/[^/]+)?$`)
	if re.MatchString(inputPath) {
		fmt.Println(aurora.Faint(("Using GitHub CLI to resolve URL...")))

		ghCloneCmd := exec.Command("gh", "repo", "view", inputPath, "--json", "url", "-q", ".url")
		out, err := ghCloneCmd.CombinedOutput()
		if err != nil {
			fmt.Println(aurora.Faint(aurora.Red("Failed to resolve repository with GitHub CLI")))
		} else {
			resolvedPath = strings.TrimSpace(string(out))
			fmt.Println(aurora.Faint(aurora.Green("GitHub URL: " + resolvedPath)))
		}
	}

	gitCloneArgs := []string{"clone"}
	if resolvedPath != "" {
		gitCloneArgs = append(gitCloneArgs, resolvedPath)
	} else {
		gitCloneArgs = append(gitCloneArgs, os.Args[1:]...)
	}

	gitCmd := exec.Command("git", gitCloneArgs...)
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	gitCmd.Run()
}
