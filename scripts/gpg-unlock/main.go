package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	helpers.EnsureAdminExecution()
	fmt.Println("Ensured admin privileges")

	out, err := exec.Command("ps", "aux").Output()
	if err == nil {
		fmt.Println("Listed processes, scanning for gpg/keyboxd")
		for _, kw := range []string{"gpg", "keyboxd"} {
			for _, line := range strings.Split(string(out), "\n") {
				if line == "" || !strings.Contains(line, kw) || strings.Contains(line, "grep") {
					continue
				}
				for _, f := range strings.Fields(line) {
					digits := true
					for i := 0; i < len(f); i++ {
						if f[i] < '0' || f[i] > '9' {
							digits = false
							break
						}
					}
					if digits {
						fmt.Println("Killing PID", f, "for", kw)
						exec.Command("sudo", "kill", "-9", f).Run()
						break
					}
				}
			}
		}
	}

	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
		if home == "" {
			home = "."
		}
	}
	fmt.Println("Using gnupg dir at", filepath.Join(home, ".gnupg"))

	entries, err := os.ReadDir(filepath.Join(home, ".gnupg"))
	if err == nil {
		for _, e := range entries {
			if e.Type().IsRegular() && filepath.Ext(e.Name()) == ".lock" {
				fmt.Println("Removing lock file", e.Name())
				os.Remove(filepath.Join(home, ".gnupg", e.Name()))
			}
		}
	}
}
