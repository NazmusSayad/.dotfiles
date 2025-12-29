package main

import (
	helpers "dotfiles/src/helpers"
	"fmt"
	"os"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func main() {
	helpers.EnsureAdminExecution()
	symlinkConfigs := helpers.ReadConfig[[]SymlinkConfig]("@/config/symlink.jsonc")

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)
		targetPath := helpers.ResolvePath(config.Target)

		fmt.Println()
		helpers.GenerateSymlink(sourcePath, targetPath)
	}
}
