package main

import (
	"dotfiles/src/helpers"
	"os"
	"os/exec"
)

func main() {
	executable := os.Args[1]
	arguments := os.Args[2:]

	cmd := exec.Command(executable, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()

	helpers.PressAnyKeyOrWaitToExit()
}
