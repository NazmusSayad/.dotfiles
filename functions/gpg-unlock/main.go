package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	psOut, _ := exec.Command("ps", "aux").Output()
	for _, line := range strings.Split(string(psOut), "\n") {
		if !strings.Contains(line, "gpg") || strings.Contains(line, "grep") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		pid := fields[0]
		fmt.Println("Found GPG process with PID: " + pid)
		exec.Command("kill", "-9", pid).Run()
	}

	for _, line := range strings.Split(string(psOut), "\n") {
		if !strings.Contains(line, "keyboxd") || strings.Contains(line, "grep") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		pid := fields[0]
		fmt.Println("Found keyboxd process with PID: " + pid)
		exec.Command("kill", "-9", pid).Run()
	}

	home := os.Getenv("HOME")
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	if home != "" {
		lockfiles, _ := filepath.Glob(filepath.Join(home, ".gnupg", "*.lock"))
		for _, lf := range lockfiles {
			os.Remove(lf)
		}
	}
}

