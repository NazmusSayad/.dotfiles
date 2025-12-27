package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	aliasCommand := "{COMMAND}"
	aliasArguments := []string{"{ARGUMENTS}"}
	scriptArguments := os.Args[1:]

	println(aliasCommand)
	println(aliasArguments)
	println(scriptArguments)

	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: aliasCommand,
		Args:    append(aliasArguments, scriptArguments...),
		Exit:    true,
	})
}
