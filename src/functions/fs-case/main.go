package main

import "dotfiles/src/helpers"

func main() {
	helpers.ExecWithNativeOutputAndExit("fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive")
}
