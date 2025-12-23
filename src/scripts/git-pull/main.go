package main

import (
	"dotfiles/src/helpers"
)

func main() {
	helpers.SimulateCommandAlias([]string{"git", "pull"})
}
