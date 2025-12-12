package main

import (
	helpers "dotfiles/src/helpers"
	"os"
)

func main() {
	arguments := os.Args[1:]
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "docker",
		Args:    append([]string{"compose"}, arguments...),
		Exit:    true,
	})
}
