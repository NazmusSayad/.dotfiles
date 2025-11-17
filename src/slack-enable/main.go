package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	helpers "dotfiles/src"
)

func main() {
	helpers.EnsureAdminExecution()
	cmd := exec.Command("schtasks", "/Change", "/TN", "\\task-slack-start", "/Enable")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	script := filepath.Join(cwd, "./slack/start.vbs")
	start := exec.Command("cscript", "//nologo", script)
	start.Stdout = os.Stdout
	start.Stderr = os.Stderr

	if err := start.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "cscript failed:", err)
		helpers.PressAnyKeyOrWaitToExit()
	}
}
