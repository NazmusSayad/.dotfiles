package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	aliasCommand := ""
	aliasArguments := []string{}

	scriptArguments := os.Args[1:]

	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: aliasCommand,
		Args:    append(aliasArguments, scriptArguments...),
		Exit:    true,
	})
}
