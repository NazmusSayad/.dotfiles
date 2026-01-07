package main

import (
	"dotfiles/src/helpers/scoop"
	"fmt"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	configs := scoop.ReadScoopConfig()
	exports := scoop.GetScoopExports()

	fmt.Println()
	fmt.Println(aurora.Green("Syncing buckets..."))
	scoop.SyncBuckets(configs, exports)

	fmt.Println()
	fmt.Println(aurora.Green("Syncing apps..."))
	scoop.SyncApps(configs, exports)

	fmt.Println()
	fmt.Println(aurora.Green("✅ Scoop installed successfully"))
}
