package main

import (
	"fmt"
	"os"

	"dotfiles/src/helpers"
)

func main() {
	if err := helpers.PowerRestart(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
