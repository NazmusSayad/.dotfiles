package main

import (
	"os"

	"dotfiles/src/helpers"

	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "git-back <commit>",
		Short: "Restore working tree files from a commit",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			helpers.ExecNativeCommand(
				[]string{"git", "restore", "--source", args[0], "--", "."},
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
