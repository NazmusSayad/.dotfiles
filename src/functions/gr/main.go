package main

import (
	"bufio"
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	fmt.Println(aurora.Red("Restore and clean?"))

	fmt.Println("Press [Enter] to confirm, or any other key to cancel: ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if strings.TrimRight(line, "\r\n") != "" {
		fmt.Println(aurora.Red("Aborted."))
		os.Exit(0)
	}

	helpers.ExecWithNativeOutput("git", "restore", ".")
	helpers.ExecWithNativeOutputAndExit("git", "clean", "-fd")
}
