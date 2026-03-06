package main

import (
	"os"

	"dotfiles/src/helpers"
)

func main() {
	helpers.ExecNativeCommand(os.Args[1:])
	helpers.PressAnyKeyOrWaitToExit()
}
