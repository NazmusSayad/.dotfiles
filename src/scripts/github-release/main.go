package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"dotfiles/src/helpers"

	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: github-release [release]")
		os.Exit(1)
	}

	helpers.InGitRepoOrExit()

	reader := bufio.NewReader(os.Stdin)
	tagFromArgs := len(os.Args) == 2
	tag := ""
	if tagFromArgs {
		tag = strings.TrimSpace(os.Args[1])
	} else {
		releasesOutput, err := exec.Command(
			"gh",
			"release",
			"list",
			"--limit",
			"5",
			"--json",
			"tagName,publishedAt",
		).Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to list recent releases"))
			os.Exit(1)
		}

		var releases []struct {
			TagName     string    `json:"tagName"`
			PublishedAt time.Time `json:"publishedAt"`
		}
		if err := json.Unmarshal(releasesOutput, &releases); err != nil {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to read recent releases"))
			os.Exit(1)
		}
		sort.Slice(releases, func(i, j int) bool {
			return releases[i].PublishedAt.Before(releases[j].PublishedAt)
		})

		if len(releases) == 0 {
			fmt.Println(aurora.Faint("No existing releases"))
		} else {
			fmt.Println(aurora.Blue("Recent:").Bold())

			longestTag := 0
			for _, release := range releases {
				if len(release.TagName) > longestTag {
					longestTag = len(release.TagName)
				}
			}

			for _, release := range releases {
				paddedTag := release.TagName + strings.Repeat(" ", longestTag-len(release.TagName))
				fmt.Println(
					aurora.Cyan(paddedTag),
					aurora.Faint(humanize.Time(release.PublishedAt)),
				)
			}
		}

		fmt.Println()
		fmt.Print(aurora.Green("Name: "))
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to read release"))
			os.Exit(1)
		}
		tag = strings.TrimSpace(input)
	}

	if tag == "" {
		fmt.Fprintln(os.Stderr, aurora.Red("Release cannot be empty"))
		os.Exit(1)
	}

	if exec.Command("gh", "release", "view", tag, "--json", "tagName").Run() == nil {
		fmt.Print(aurora.Yellow("Release " + tag + " already exists. Recreate it? [y/N]: ").Bold())
		confirmation, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
			os.Exit(1)
		}

		if !strings.EqualFold(strings.TrimSpace(confirmation), "y") {
			fmt.Println(aurora.Red("Release recreation cancelled"))
			return
		}

		helpers.ExecNativeCommand(
			[]string{"gh", "release", "delete", tag, "--yes"},
			helpers.ExecCommandOptions{Exit: true},
		)

		helpers.ExecNativeCommand(
			[]string{
				"gh",
				"api",
				"--method",
				"DELETE",
				"repos/{owner}/{repo}/git/refs/tags/" + url.PathEscape(tag),
			},
			helpers.ExecCommandOptions{Exit: true},
		)
	} else if tagFromArgs {
		fmt.Print(
			"Create release ",
			aurora.Cyan(tag).Bold(),
			" ",
			aurora.Faint("[Enter]: "),
		)
		confirmation, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
			os.Exit(1)
		}

		if strings.TrimRight(confirmation, "\r\n") != "" {
			fmt.Println(aurora.Red("Release creation cancelled"))
			return
		}
	}

	defaultBranchOutput, err := exec.Command(
		"gh",
		"repo",
		"view",
		"--json",
		"defaultBranchRef",
		"--jq",
		".defaultBranchRef.name",
	).Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to resolve default branch"))
		os.Exit(1)
	}

	defaultBranch := strings.TrimSpace(string(defaultBranchOutput))
	if defaultBranch == "" {
		fmt.Fprintln(os.Stderr, aurora.Red("Default branch not found"))
		os.Exit(1)
	}

	helpers.ExecNativeCommand(
		[]string{"gh", "release", "create", tag, "--target", defaultBranch, "--generate-notes"},
		helpers.ExecCommandOptions{Exit: true},
	)
}
