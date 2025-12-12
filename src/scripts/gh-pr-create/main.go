package main

import (
	helpers "dotfiles/src/helpers"
	"os"
)

func main() {
	arguments := os.Args[1:]
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "gh",
		Args:    append([]string{"pr", "create", "-B"}, arguments...),
		Exit:    true,
	})
}
