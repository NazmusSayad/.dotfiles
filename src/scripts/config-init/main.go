package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"
	"dotfiles/src/helpers/symlink"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if !helpers.IsWSL() {
		helpers.EnsureAdminExecution()
	}
	if err := helpers.ValidateInvokingUser(); err != nil {
		fmt.Println(aurora.Red("UNEXPECTED: " + err.Error()))
		os.Exit(1)
	}
	symlinkConfigs := symlink.ReadConfigs()

	if len(symlinkConfigs) == 0 {
		fmt.Println("No symlink configurations found.")
		os.Exit(1)
	}

	newlyCreatedFiles := []string{}

	for _, config := range symlinkConfigs {
		sourcePath := helpers.ResolvePath(config.Source)
		fmt.Println(aurora.Blue(config.Source))

		for _, target := range config.LinkTargets {
			targetPath := helpers.ResolvePath(target)
			if helpers.GenerateSymlink(sourcePath, targetPath, config.InheritPerm) == nil {
				newlyCreatedFiles = append(newlyCreatedFiles, targetPath)
				fmt.Println(aurora.Blue("->"), aurora.Faint(target))
			}
		}

		for _, target := range config.CopyTargets {
			targetPath := helpers.ResolvePath(target)
			if helpers.CopyFile(sourcePath, targetPath, config.InheritPerm) == nil {
				newlyCreatedFiles = append(newlyCreatedFiles, targetPath)
				fmt.Println(aurora.Blue("=>"), aurora.Faint(targetPath))
			}
		}

		fmt.Println()
	}

	if !helpers.IsWSL() {
		lockFile, _ := os.ReadFile(helpers.ResolvePath("@/.local/symlink.lock"))
		for _, file := range strings.Split(string(lockFile), "\n") {
			if file != "" && !slices.Contains(newlyCreatedFiles, file) {
				fmt.Println(aurora.Yellow("Deleting stale link: " + file))
				os.RemoveAll(file)
			}
		}

		lockPath := helpers.ResolvePath("@/.local/symlink.lock")
		os.WriteFile(lockPath, []byte(strings.Join(newlyCreatedFiles, "\n")), 0o644)
		helpers.ApplyUserOwnership(lockPath)
	}

	fmt.Println()
	opencode.Configure()
	configureWSL()
}

func configureWSL() {
	if runtime.GOOS != "windows" {
		return
	}

	fmt.Println()
	fmt.Println(aurora.Blue("Updating WSL configuration"))
	cmd := exec.Command(
		"wsl.exe",
		"--distribution", "Ubuntu-26.04",
		"--user", "sayad",
		"--cd", "~",
		"--exec", "/bin/bash", "-lc",
		`cd "$HOME/.dotfiles" && "$HOME/.local/bin/mise" exec -- go run ./src/scripts/config-init/main.go`,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(aurora.Red("Failed to update WSL configuration: " + err.Error()))
	}
}
