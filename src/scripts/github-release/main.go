package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: github-release [tag]")
		os.Exit(1)
	}

	helpers.InGitRepoOrExit()

	tag := ""
	if len(os.Args) == 2 {
		tag = strings.TrimSpace(os.Args[1])
	} else {
		tagsOutput, err := exec.Command("git", "tag", "--sort=-creatordate").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to list tags"))
			os.Exit(1)
		}

		tags := strings.Fields(string(tagsOutput))
		if len(tags) > 5 {
			tags = tags[:5]
		}

		if len(tags) == 0 {
			fmt.Println(aurora.Faint("No existing tags"))
		} else {
			fmt.Println(aurora.Faint("Recent tags:"))
			for _, recentTag := range tags {
				fmt.Println("  " + recentTag)
			}
		}

		fmt.Print("New tag: ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to read tag"))
			os.Exit(1)
		}
		tag = strings.TrimSpace(input)
	}

	if tag == "" {
		fmt.Fprintln(os.Stderr, aurora.Red("Tag cannot be empty"))
		os.Exit(1)
	}

	helpers.ExecNativeCommand(
		[]string{"gh", "release", "create", tag, "--generate-notes"},
		helpers.ExecCommandOptions{Exit: true},
	)
}
