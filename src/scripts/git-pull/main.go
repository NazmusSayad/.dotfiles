package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "git",
		Args:    append([]string{"pull"}, os.Args...),
		Exit:    true,
	})
}
