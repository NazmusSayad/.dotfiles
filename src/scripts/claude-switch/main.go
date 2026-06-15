package main

import (
	"fmt"
	"os"

	"dotfiles/src/helpers/claude"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	localDir := claude.ResolveLocalDir()
	targetPath := claude.ResolveCredentialsPath()

	accounts := claude.ReadAccounts(localDir, targetPath)
	if len(accounts) == 0 {
		fmt.Println(aurora.Red("No *.credentials.json files found in .local"))
		return
	}

	fmt.Println("> Current Claude account: " + aurora.Green(claude.CurrentAccount(accounts)).String())

	last := claude.ReadLastAccount(localDir)
	if last == "" {
		last = "unknown"
	}
	fmt.Println("> Last used account: " + aurora.Cyan(last).String())
	fmt.Println()

	choice := claude.SelectAccount(accounts)
	if choice == nil {
		return
	}

	data, err := os.ReadFile(choice.Path)
	if err != nil {
		fmt.Println(aurora.Red("Failed to read source: " + choice.Path))
		return
	}

	err = os.WriteFile(targetPath, data, 0o600)
	if err != nil {
		fmt.Println(aurora.Red("Failed to write target: " + targetPath))
		return
	}
	fmt.Println("> Switched Claude account to " + aurora.Green(choice.Name).String())

	claude.WriteLastAccount(localDir, choice.Name)
}
