package main

import (
	"dotfiles/src/helpers"
)

func main() {
	helpers.ExecNativeCommand([]string{
		"powershell", "-NoProfile", "-Command",
		"Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser; Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression",
	}, helpers.ExecCommandOptions{
		Exit: true,
	})
}
