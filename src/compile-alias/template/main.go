package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	aliasCommand := "{COMMAND}"
	scriptArguments := append([]string{"{ARGUMENTS}"}, os.Args[1:]...)

	helpers.ExecNativeCommand(
		append([]string{aliasCommand}, scriptArguments...),
		helpers.ExecCommandOptions{
			Exit: true,
		},
	)
}
