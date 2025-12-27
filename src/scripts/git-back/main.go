package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	commitHash := ""
	if len(os.Args) > 1 {
		commitHash = os.Args[1]
	}

	if commitHash == "" {
		fmt.Println(aurora.Red("Commit hash required"))
		os.Exit(1)
	}

	helpers.ExecNativeCommand(
		[]string{"git", "restore", "--source", commitHash, "--", "."},
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}
