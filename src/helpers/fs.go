package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

func CopyFile(source string, target string) error {
	if !utils.IsFileExists(source) {
		fmt.Println(aurora.Red("UNEXPECTED: Source not found: " + source))
		return fmt.Errorf("source not found: %s", source)
	}

	if utils.IsFileExists(target) {
		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println(aurora.Red("UNEXPECTED: Error deleting target: " + target))
			return removeErr
		}
	}

	createdDirs, ok := createTargetDirs(target)
	if !ok {
		return fmt.Errorf("failed to create target directories for: %s", target)
	}

	content, err := os.ReadFile(source)
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error reading source: " + err.Error()))
		return fmt.Errorf("error reading source: %w", err)
	}

	info, err := os.Stat(source)
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error stating source: " + err.Error()))
		return fmt.Errorf("error stating source: %w", err)
	}

	err = os.WriteFile(target, content, info.Mode().Perm())
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error copying file: " + err.Error()))
		return fmt.Errorf("error copying file: %w", err)
	}

	inheritOwnership(target, createdDirs)

	fmt.Println(aurora.Blue(source), aurora.Green("=>"), aurora.Cyan(target))
	return nil
}

func GenerateSymlink(source string, target string) error {
	if !utils.IsFileExists(source) {
		fmt.Println(aurora.Red("UNEXPECTED: Source not found: " + source))
		return fmt.Errorf("source not found: %s", source)
	}

	if utils.IsFileExists(target) {
		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println(aurora.Red("UNEXPECTED: Error deleting target: " + target))
			return removeErr
		}
	}

	createdDirs, ok := createTargetDirs(target)
	if !ok {
		return fmt.Errorf("failed to create target directories for: %s", target)
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error creating symlink: " + err.Error()))
		return fmt.Errorf("error creating symlink: %w", err)
	}

	inheritOwnership(target, createdDirs)

	fmt.Println(aurora.Blue(source), aurora.Green("->"), aurora.Cyan(target))
	return nil
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
