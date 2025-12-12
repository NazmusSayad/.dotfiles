package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println(aurora.Red("Restore and clean?"))

	fmt.Print("Press [Enter] to confirm, or any other key to cancel: ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimRight(line, "\r\n") != "" {
		fmt.Println(aurora.Red("❌ Aborted."))
		os.Exit(0)
	}

	exec.Command("git", "restore", ".").Run()
	cmd := exec.Command("git", "clean", "-fd")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		}
		os.Exit(1)
	}
}
