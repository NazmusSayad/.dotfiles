package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/helpers/symlink"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	helpers.EnsureAdminExecution()
	symlinkConfigs := symlink.ReadConfigs()

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	newlyCreatedFiles := []string{}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)

		for _, target := range config.LinkTargets {
			targetPath := helpers.ResolvePath(target)
			if helpers.GenerateSymlink(sourcePath, targetPath) == nil {
				newlyCreatedFiles = append(newlyCreatedFiles, targetPath)
			}
		}

		for _, target := range config.CopyTargets {
			targetPath := helpers.ResolvePath(target)
			if helpers.CopyFile(sourcePath, targetPath) == nil {
				newlyCreatedFiles = append(newlyCreatedFiles, targetPath)
			}
		}
	}

	lockFile, _ := os.ReadFile(helpers.ResolvePath("@/.local/symlink.lock"))
	for _, file := range strings.Split(string(lockFile), "\n") {
		if file != "" && !slices.Contains(newlyCreatedFiles, file) {
			fmt.Println(aurora.Yellow("Deleting stale link: " + file))
			os.RemoveAll(file)
		}
	}

	os.WriteFile(
		helpers.ResolvePath("@/.local/symlink.lock"),
		[]byte(strings.Join(newlyCreatedFiles, "\n")),
		0o644,
	)
}
