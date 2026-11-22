package main

import (
	"os"

	"dotfiles/src/helpers/agentstore"
)

func main() {
	agentstore.Push(agentstore.LoadConfig(os.Args[1:]))
}
