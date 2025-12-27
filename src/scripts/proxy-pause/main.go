package main

import (
	"dotfiles/src/helpers"
	"os"
)

func main() {
	helpers.ExecNativeCommand(os.Args[1:])
	helpers.PressAnyKeyOrWaitToExit()
}
