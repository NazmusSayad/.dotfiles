package main

import (
	"os"

	"dotfiles/src/helpers/agentstore"
)

func main() {
	agentstore.Pull(agentstore.LoadConfig(os.Args[1:]))
}
