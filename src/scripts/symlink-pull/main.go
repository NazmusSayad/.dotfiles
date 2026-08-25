package main

import (
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/helpers/symlink"
	"dotfiles/src/utils"
)

func main() {
	helpers.EnsureAdminExecution()
	symlinkConfigs := symlink.ReadConfigs()

	if len(symlinkConfigs) == 0 {
		fmt.Println("No copy configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		if len(config.CopyTargets) == 0 {
			continue
		}

		sourcePath := helpers.ResolvePath(config.Source)
		targetPath := helpers.ResolvePath(config.CopyTargets[0])

		if !utils.IsFileExists(targetPath) {
			fmt.Println("Skipping, target not found:", targetPath)
			continue
		}

		helpers.CopyFile(targetPath, sourcePath)
	}
}
