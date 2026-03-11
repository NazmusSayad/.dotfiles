package main

import (
	"fmt"
	"os"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	symLinkErr := helpers.ExecNativeCommand([]string{"symlink-setup"})
	if symLinkErr != nil {
		fmt.Println(aurora.Red("Error setting up symbolic links:"), aurora.Bold(symLinkErr.Error()))
		os.Exit(1)
	}

	codeExtErr := helpers.ExecNativeCommand([]string{"code-ext-sync"})
	if codeExtErr != nil {
		fmt.Println(aurora.Red("Error syncing code extensions:"), aurora.Bold(codeExtErr.Error()))
		os.Exit(1)
	}

	codeStateErr := helpers.ExecNativeCommand([]string{"code-state-push"})
	if codeStateErr != nil {
		fmt.Println(aurora.Red("Error pushing code state:"), aurora.Bold(codeStateErr.Error()))
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Code environment initialized successfully!"))
}
