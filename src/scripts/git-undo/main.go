package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	undoCount := 1

	if len(os.Args) > 1 {
		count, err := strconv.Atoi(os.Args[1])
		if err != nil || count <= 0 {
			fmt.Println(aurora.Red("Invalid undo count"))
			os.Exit(1)
		}

		undoCount = count
	}

	commitCount := aurora.Bold(strconv.Itoa(undoCount)).String()
	commitWord := utils.Ternary(undoCount > 1, "commits", "commit")
	fmt.Println(aurora.Red("Undoing last " + commitCount + " " + commitWord + "..."))

	helpers.ExecNativeCommand(
		[]string{"git", "reset", "--soft", "HEAD~" + strconv.Itoa(undoCount)},
		helpers.ExecCommandOptions{Exit: true},
	)
}
