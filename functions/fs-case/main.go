package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}
