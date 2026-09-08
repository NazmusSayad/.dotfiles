package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	"github.com/otiai10/copy"
)

func CopyFile(source string, target string, inheritPerm bool) error {
	if !utils.IsFileExists(source) {
		fmt.Println(aurora.Red("UNEXPECTED: Source not found: " + source))
		return fmt.Errorf("source not found: %s", source)
	}
	if !inheritPerm {
		if err := ValidateInvokingUser(); err != nil {
			fmt.Println(aurora.Red("UNEXPECTED: " + err.Error()))
			return err
		}
	}

	if utils.IsFileExists(target) {
		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println(aurora.Red("UNEXPECTED: Error deleting target: " + target))
			return removeErr
		}
	}

	createdDirs, ok := createTargetDirs(target)
	if !inheritPerm {
		defer ApplyUserOwnership(append(createdDirs, target)...)
	}
	if !ok {
		return fmt.Errorf("failed to create target directories for: %s", target)
	}

	err := copy.Copy(source, target, copy.Options{
		PreserveTimes: true,
		PreserveOwner: inheritPerm,
	})
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error copying source: " + err.Error()))
		return fmt.Errorf("error copying source: %w", err)
	}

	if inheritPerm {
		inheritDirOwnership(createdDirs)
	}

	fmt.Println(aurora.Blue(source), aurora.Green("=>"), aurora.Cyan(target))
	return nil
}

func GenerateSymlink(source string, target string, inheritPerm bool) error {
	if !utils.IsFileExists(source) {
		fmt.Println(aurora.Red("UNEXPECTED: Source not found: " + source))
		return fmt.Errorf("source not found: %s", source)
	}
	if !inheritPerm {
		if err := ValidateInvokingUser(); err != nil {
			fmt.Println(aurora.Red("UNEXPECTED: " + err.Error()))
			return err
		}
	}

	if utils.IsFileExists(target) {
		removeErr := os.RemoveAll(target)
		if removeErr != nil {
			fmt.Println(aurora.Red("UNEXPECTED: Error deleting target: " + target))
			return removeErr
		}
	}

	createdDirs, dirsOk := createTargetDirs(target)
	if !inheritPerm {
		defer ApplyUserOwnership(append(createdDirs, target)...)
	}
	if !dirsOk {
		return fmt.Errorf("failed to create target directories for: %s", target)
	}

	err := os.Symlink(source, target)
	if err != nil {
		fmt.Println(aurora.BrightRed("UNEXPECTED: Error creating symlink: " + err.Error()))
		return fmt.Errorf("error creating symlink: %w", err)
	}

	if inheritPerm {
		inheritOwnership(target, createdDirs)
	}

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
		return missing, false
	}

	for i, j := 0, len(missing)-1; i < j; i, j = i+1, j-1 {
		missing[i], missing[j] = missing[j], missing[i]
	}

	return missing, true
}
