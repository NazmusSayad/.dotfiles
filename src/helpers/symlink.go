package helpers

import (
	"dotfiles/src/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/logrusorgru/aurora/v4"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func GenerateSymlink(source string, target string) {
	fmt.Println(aurora.Faint("Symlinking: " + source))

	if !utils.IsFileExists(source) {
		fmt.Println(aurora.Red("UNEXPECTED: Source not found: " + source))
		return
	}

	if utils.IsFileExists(target) {
		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println(aurora.Red("UNEXPECTED: Error deleting target: " + target))
			return
		}
	}

	targetDir := filepath.Dir(target)
	if !utils.IsFileExists(targetDir) {
		mkdirErr := os.MkdirAll(targetDir, 0755)
		if mkdirErr != nil {
			fmt.Println(aurora.Red("UNEXPECTED: Error creating target directory: " + targetDir))
			return
		}
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error creating symlink: " + err.Error()))
		return
	}

	fmt.Println(aurora.Green("-> " + target))
}
