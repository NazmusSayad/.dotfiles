package main

import (
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/helpers/symlink"
)

func main() {
	helpers.EnsureAdminExecution()
	symlinkConfigs := symlink.ReadConfigs()

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)

		for _, target := range config.Targets {
			targetPath := helpers.ResolvePath(target)

			if config.Copy {
				helpers.CopyFile(sourcePath, targetPath)
			} else {
				helpers.GenerateSymlink(sourcePath, targetPath)
			}
		}
	}
}
