package main

import (
	helpers "dotfiles/src"
)

func main() {
	helpers.EnsureAdminExecution()

	// TODO: Implement msys-setup

	helpers.WaitForInputAndExit()
}
