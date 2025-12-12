package helpers

import (
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
	fmt.Println(aurora.Green("Symlinking: " + source + " -> " + target))

	if !IsFileExists(source) {
		fmt.Println("UNEXPECTED: Source not found: " + source)
		return
	}

	if IsFileExists(target) {
		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println("UNEXPECTED: Error deleting target: " + target)
			return
		}
	}

	targetDir := filepath.Dir(target)
	if !IsFileExists(targetDir) {
		mkdirErr := os.MkdirAll(targetDir, 0755)
		if mkdirErr != nil {
			fmt.Println("UNEXPECTED: Error creating target directory: " + targetDir)
			return
		}
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println("UNEXPECTED: Error creating symlink: " + err.Error())
		return
	}

	fmt.Println(source, "->", target)
}
