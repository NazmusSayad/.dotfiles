package main

import (
	helpers "dotfiles/src/helpers"
)

func main() {
	helpers.SimulateCommandAlias([]string{"gh", "repo", "view", "--web"})
}
