package main

import (
	"dotfiles/src/helpers/scoop"
)

func main() {
	configs := scoop.ReadScoopConfig()
	exports := scoop.GetScoopExports()

	scoop.SyncBuckets(configs, exports)
	scoop.SyncApps(configs, exports)
}
