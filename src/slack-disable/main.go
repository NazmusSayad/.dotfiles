package main

import (
	"os"
	"os/exec"

	helpers "dotfiles/src"
)

func main() {
	helpers.EnsureAdminExecution()
	cmd := exec.Command("schtasks", "/Change", "/TN", "\\task-slack-start", "/Disable")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
