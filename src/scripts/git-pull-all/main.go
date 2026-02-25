package main

import (
	"dotfiles/src/helpers"
	"fmt"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println(aurora.Yellow("Pulling changes from all branches"))

	helpers.ExecNativeCommand(
		[]string{"git", "pull", "--all"},
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}
