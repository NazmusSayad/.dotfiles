package main

import (
	helpers "dotfiles/src/helpers"
	"os"
)

func main() {
	arguments := os.Args[1:]
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "winget",
		Args:    append([]string{"alias"}, arguments...),
		Exit:    true,
	})
}
