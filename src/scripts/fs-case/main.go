package main

import "dotfiles/src/helpers"

func main() {
	helpers.ExecNativeCommand(helpers.ExecCommandOptions{
		Command: "fsutil.exe",
		Args:    []string{"file", "setCaseSensitiveInfo", ".", "enable", "recursive"},
		Exit:    true,
	})
}
