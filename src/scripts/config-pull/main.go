package main

import (
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/helpers/symlink"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
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

		targetRaw := config.CopyTargets[0]
		if targetRaw == "" {
			fmt.Println("Skipping, empty target for source:", config.Source)
			continue
		}

		targetPath := helpers.ResolvePath(targetRaw)
		sourcePath := helpers.ResolvePath(config.Source)

		if !utils.IsFileExists(targetPath) {
			fmt.Println("Skipping, target not found:", targetPath)
			continue
		}

		helpers.CopyFile(targetPath, sourcePath, config.InheritPerm)

		fmt.Println(aurora.Blue(targetPath))
		fmt.Println(aurora.Blue("=>"), aurora.Faint(config.Source))
		fmt.Println()
	}
}
