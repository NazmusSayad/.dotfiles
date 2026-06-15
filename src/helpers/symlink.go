package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/src/utils"

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

	createdDirs, ok := createTargetDirs(target)
	if !ok {
		return
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error creating symlink: " + err.Error()))
		return
	}

	inheritOwnership(target, createdDirs)

	fmt.Println(aurora.Green("-> " + target))
}

func createTargetDirs(target string) ([]string, bool) {
	targetDir := filepath.Dir(target)
	if utils.IsFileExists(targetDir) {
		return nil, true
	}

	var missing []string
	for dir := targetDir; !utils.IsFileExists(dir); dir = filepath.Dir(dir) {
		missing = append(missing, dir)
		if filepath.Dir(dir) == dir {
			break
		}
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		fmt.Println(aurora.Red("UNEXPECTED: Error creating target directory: " + targetDir))
		return nil, false
	}

	for i, j := 0, len(missing)-1; i < j; i, j = i+1, j-1 {
		missing[i], missing[j] = missing[j], missing[i]
	}

	return missing, true
}
