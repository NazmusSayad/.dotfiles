package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: os.Args[1],
		Args:    os.Args[2:],
	})

	helpers.PressAnyKeyOrWaitToExit()
}
