package main

import "dotfiles/src/helpers"

func main() {
	helpers.SimulateCommandAlias([]string{"fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive"})
}
