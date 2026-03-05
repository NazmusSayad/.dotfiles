package main

import (
	"dotfiles/src/helpers"
)

func main() {
	helpers.ExecNativeCommand(
		[]string{"net", "stop", "winnat"}, helpers.ExecCommandOptions{
			Exit:    true,
			AsAdmin: true,
		},
	)
}
