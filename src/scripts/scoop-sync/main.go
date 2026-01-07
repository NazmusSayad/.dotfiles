package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/helpers/scoop"
	"fmt"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	configs := scoop.ReadScoopAppConfig()
	exports := scoop.GetScoopExports()

	fmt.Println()
	fmt.Println(aurora.Green("Syncing buckets..."))
	scoop.SyncBuckets(configs, exports)

	fmt.Println()
	fmt.Println(aurora.Green("Syncing apps..."))
	scoop.SyncApps(configs, exports)

	fmt.Println()
	fmt.Println(aurora.Green("Updating scoop..."))
	helpers.ExecNativeCommand([]string{"scoop", "update"})

	fmt.Println()
	fmt.Println(aurora.Green("Updating apps..."))
	helpers.ExecNativeCommand([]string{
		"scoop", "update", "*", "--no-cache",
	}, helpers.ExecCommandOptions{
		Exit: true,
	})
}
