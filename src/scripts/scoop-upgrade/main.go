package main

import "dotfiles/src/helpers"

func main() {
	helpers.ExecNativeCommand([]string{"scoop", "update"})
	helpers.ExecNativeCommand([]string{
		"scoop", "update", "*", "--no-cache",
	}, helpers.ExecCommandOptions{
		Exit: true,
	})
}
