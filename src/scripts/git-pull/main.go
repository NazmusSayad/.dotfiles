package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	arguments := os.Args[1:]
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "git",
		Args:    append([]string{"pull"}, arguments...),
		Exit:    true,
	})
}
