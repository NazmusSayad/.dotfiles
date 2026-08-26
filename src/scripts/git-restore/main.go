package main

import (
	"errors"
	"fmt"
	"os"

	"dotfiles/src/helpers"

	"charm.land/huh/v2"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	confirmed := true
	err := huh.NewConfirm().
		Title("Restore and clean? ").
		Inline(true).
		Value(&confirmed).
		WithTheme(helpers.HuhTheme()).
		Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println(aurora.Red("Aborted."))
			return
		}
		fmt.Fprintln(os.Stderr, aurora.Red("Failed to read confirmation"))
		os.Exit(1)
	}
	if !confirmed {
		fmt.Println(aurora.Red("Aborted."))
		return
	}

	helpers.ExecNativeCommand([]string{"git", "restore", "."})
	helpers.ExecNativeCommand(
		[]string{"git", "clean", "-fd"},
		helpers.ExecCommandOptions{Exit: true},
	)
}
