package main

import (
	helpers "dotfiles/src/helpers"
)

func main() {
	helpers.SimulateCommandAlias([]string{"gh", "pr", "create", "-B"})
}
