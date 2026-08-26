package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"dotfiles/src/helpers"

	"charm.land/huh/v2"
	gh "github.com/cli/go-gh/v2"
	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "github-release [release]",
		Short: "Create a GitHub release",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			helpers.InGitRepoOrExit()

			ctx := context.Background()
			tagFromArgs := len(args) == 1
			tag := ""
			if tagFromArgs {
				tag = strings.TrimSpace(args[0])
			} else {
				releasesOutput, _, err := gh.Exec(
					"release",
					"list",
					"--limit",
					"5",
					"--json",
					"tagName,publishedAt",
				)
				if err != nil {
					fmt.Fprintln(os.Stderr, aurora.Red("Failed to list recent releases"))
					os.Exit(1)
				}

				var releases []struct {
					TagName     string    `json:"tagName"`
					PublishedAt time.Time `json:"publishedAt"`
				}
				if err := json.Unmarshal(releasesOutput.Bytes(), &releases); err != nil {
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
				err = huh.NewInput().
					Title(aurora.Green("Name").String()).
					Inline(true).
					Value(&tag).
					WithTheme(helpers.HuhTheme()).
					Run()
				if err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						fmt.Println(aurora.Red("Release creation cancelled"))
						return
					}
					fmt.Fprintln(os.Stderr, aurora.Red("Failed to read release"))
					os.Exit(1)
				}
				tag = strings.TrimSpace(tag)
			}

			if tag == "" {
				fmt.Fprintln(os.Stderr, aurora.Red("Release cannot be empty"))
				os.Exit(1)
			}

			if _, _, err := gh.Exec("release", "view", tag, "--json", "tagName"); err == nil {
				recreate := false
				err := huh.NewConfirm().
					Title(aurora.Yellow("Release "+tag+" already exists. Recreate it?").Bold().String() + " ").
					Inline(true).
					Value(&recreate).
					WithTheme(helpers.HuhTheme()).
					Run()
				if err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						fmt.Println(aurora.Red("Release recreation cancelled"))
						return
					}
					fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
					os.Exit(1)
				}

				if !recreate {
					fmt.Println(aurora.Red("Release recreation cancelled"))
					return
				}

				if err := gh.ExecInteractive(ctx, "release", "delete", tag, "--yes"); err != nil {
					os.Exit(1)
				}

				if err := gh.ExecInteractive(
					ctx,
					"api",
					"--method",
					"DELETE",
					"repos/{owner}/{repo}/git/refs/tags/"+url.PathEscape(tag),
				); err != nil {
					os.Exit(1)
				}
			} else if tagFromArgs {
				confirmed := true
				err := huh.NewConfirm().
					Title(fmt.Sprint("Create release ", aurora.Cyan(tag).Bold(), " ")).
					Inline(true).
					Value(&confirmed).
					WithTheme(helpers.HuhTheme()).
					Run()
				if err != nil {
					if errors.Is(err, huh.ErrUserAborted) {
						fmt.Println(aurora.Red("Release creation cancelled"))
						return
					}
					fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
					os.Exit(1)
				}

				if !confirmed {
					fmt.Println(aurora.Red("Release creation cancelled"))
					return
				}
			}

			defaultBranchOutput, _, err := gh.Exec(
				"repo",
				"view",
				"--json",
				"defaultBranchRef",
				"--jq",
				".defaultBranchRef.name",
			)
			if err != nil {
				fmt.Fprintln(os.Stderr, aurora.Red("Failed to resolve default branch"))
				os.Exit(1)
			}

			defaultBranch := strings.TrimSpace(defaultBranchOutput.String())
			if defaultBranch == "" {
				fmt.Fprintln(os.Stderr, aurora.Red("Default branch not found"))
				os.Exit(1)
			}

			if err := gh.ExecInteractive(
				ctx,
				"release",
				"create",
				tag,
				"--target",
				defaultBranch,
				"--generate-notes",
			); err != nil {
				os.Exit(1)
			}
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
