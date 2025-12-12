package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	commitHash := ""
	if len(os.Args) > 1 {
		commitHash = os.Args[1]
	}

	if commitHash == "" {
		fmt.Println(aurora.Red("❌ Commit hash required"))
		os.Exit(1)
	}

	cmd := exec.Command("git", "restore", "--source", commitHash, "--", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}
