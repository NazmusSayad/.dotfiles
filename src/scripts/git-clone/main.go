package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:                "git-clone <repository-path> [additional-arguments]",
		Short:              "Clone a Git repository",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		Run: func(_ *cobra.Command, args []string) {
			inputPath := args[0]
			resolvedPath := ""

			re := regexp.MustCompile(`^[^/]+(/[^/]+)?$`)
			if re.MatchString(inputPath) {
				fmt.Println(aurora.Faint("Using GitHub CLI to resolve URL..."))

				ghCloneCmd := exec.Command("gh", "repo", "view", inputPath, "--json", "url", "-q", ".url")
				out, err := ghCloneCmd.Output()
				if err != nil {
					fmt.Println(aurora.Faint(aurora.Red("Failed to resolve repository with GitHub CLI")))
				} else {
					resolvedPath = strings.TrimSpace(string(out))
					fmt.Println(aurora.Faint(aurora.Green("GitHub URL: " + resolvedPath)))
				}
			}

			gitCloneArgs := []string{"git", "clone"}
			if resolvedPath != "" {
				gitCloneArgs = append(gitCloneArgs, resolvedPath)
			} else {
				gitCloneArgs = append(gitCloneArgs, args...)
			}

			helpers.ExecNativeCommand(
				gitCloneArgs,
				helpers.ExecCommandOptions{
					Exit: true,
				},
			)
		},
	}

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
