package main

import (
	"dotfiles/src/helpers/scoop"
	"fmt"
)

func main() {
	configs := scoop.ReadScoopAppConfig()
	for _, config := range configs {
		fmt.Println(config.ID)
	}
}
