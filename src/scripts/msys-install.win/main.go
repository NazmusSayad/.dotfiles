package main

import (
	"fmt"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	msys2Packages := helpers.GetMsysApps()

	if len(msys2Packages) > 0 {
		msys2PackagesString := []string{}
		for _, p := range msys2Packages {
			msys2PackagesString = append(msys2PackagesString, p.ID)
		}

		fmt.Println(aurora.Faint("- Installing"), aurora.Green(strings.Join(msys2PackagesString, " ")))
		helpers.ExecNativeCommand(append([]string{"pacman", "--noconfirm", "-S", "--needed"}, msys2PackagesString...))
	}
}
