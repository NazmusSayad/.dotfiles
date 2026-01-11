package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	runCommand("Winget", []string{"winget-upgrade"})

	fmt.Println()
	runCommand("Scoop", []string{"scoop", "update"})

	fmt.Println()
	runCommand("Mise", []string{"mise", "upgrade"})

	fmt.Println()
	runCommand("Pacman", []string{"pacman", "-Syu", "--noconfirm"})
}

func runCommand(name string, args []string) {
	cmd := exec.Command("cmd.exe", "/c", strings.Join(args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	fmt.Println()
	if err == nil {
		fmt.Println(aurora.Green("✅ " + name + " upgrade completed"))
	} else {
		fmt.Println(aurora.Red("❌ Error: " + err.Error()))
	}
}
