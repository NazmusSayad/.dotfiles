package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"dotfiles/src/helpers"
)

var FILES_TO_PATCH = []string{
	"bash.exe",
	"sh.exe",
}

func main() {
	helpers.EnsureAdminExecution()

	gitExe := findExecutable("git")
	if gitExe == "" {
		panic("Git executable not found in PATH")
	}

	gitBin := filepath.Join(filepath.Dir(gitExe), "..", "bin")
	gitBin, err := filepath.Abs(gitBin)
	if err != nil {
		panic(fmt.Errorf("failed to resolve git bin path: %w", err))
	}

	fmt.Println("Git found at:", gitExe)
	fmt.Println("Git bin:", gitBin)

	for _, file := range FILES_TO_PATCH {
		sourceExe := findExecutable(file)
		if sourceExe == "" {
			fmt.Printf("Skipping %s: not found in PATH\n", file)
			continue
		}

		targetPath := filepath.Join(gitBin, file)

		fmt.Printf("\nPatching %s\n", file)
		fmt.Println("Source:", sourceExe)
		fmt.Println("Target:", targetPath)

		backupIfExists(targetPath)
		helpers.GenerateSymlink(sourceExe, targetPath)
	}

	fmt.Println("\nDone.")
}

func findExecutable(name string) string {
	if path, err := exec.LookPath(name); err == nil {
		return path
	}

	return ""
}

func backupIfExists(path string) {
	backup := path + ".backup"

	if _, err := os.Stat(path); err != nil {
		return
	}

	if _, err := os.Stat(backup); err == nil {
		fmt.Printf("Backup already exists: %s\n", backup)
		return
	}

	src, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("failed to open %s: %w", path, err))
	}
	defer src.Close()

	dst, err := os.Create(backup)
	if err != nil {
		panic(fmt.Errorf("failed to create backup %s: %w", backup, err))
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		panic(fmt.Errorf("failed to copy %s -> %s: %w", path, backup, err))
	}

	fmt.Printf("Backed up: %s -> %s\n", path, backup)
}
