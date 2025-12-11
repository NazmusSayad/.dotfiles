package main

import (
	helpers "dotfiles/src"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!, after file chanhe, auto run again")
	println(helpers.IsFileExists("./readme.md"))
}
