package helpers

import (
	"fmt"
	"os"
	"path/filepath"
)

type SymlinkConfig struct {
	Source string
	Target string
}

func GenerateSymlink(source string, target string) {
	fmt.Printf("Symlinking: %s -> %s\n", source, target)

	if !IsFileExists(source) {
		fmt.Println("UNEXPECTED: Source not found:", source)
		return
	}

	if IsFileExists(target) {
		fmt.Println("EXPECTED: Target found, deleting:", target)

		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println("UNEXPECTED: Error deleting target:", target)
			return
		}
	}

	targetDir := filepath.Dir(target)
	if !IsFileExists(targetDir) {
		fmt.Println("EXPECTED: Target directory not found, creating:", targetDir)

		mkdirErr := os.MkdirAll(targetDir, 0755)
		if mkdirErr != nil {
			fmt.Println("UNEXPECTED: Error creating target directory:", targetDir)
			return
		}
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println("UNEXPECTED: Error creating symlink", err)
		return
	}

	fmt.Println(source, "->", target)
}
