package main

import (
	"dotfiles/src/helpers"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	msys2Packages := helpers.ReadConfig[[]string]("@/config/msys2-packages.jsonc")
	if len(msys2Packages) > 0 {
		fmt.Println(aurora.Faint("- Installing"), aurora.Green(strings.Join(msys2Packages, " ")))
		helpers.ExecNativeCommand(append([]string{"pacman", "--noconfirm", "-S", "--needed"}, msys2Packages...))
	}
}
