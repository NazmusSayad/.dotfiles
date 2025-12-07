package main

import (
	helpers "dotfiles/src"
	"os"
	"path/filepath"
)

func main() {
	helpers.EnsureAdminExecution()

	cwd, err := os.Getwd()
	if err != nil {
		println("Error getting current working directory:", err)
		os.Exit(1)
	}

	userProfile, err := os.UserHomeDir()
	if err != nil {
		println("Error getting user home directory:", err)
		os.Exit(1)
	}

	dotfilesDir := filepath.Join(userProfile, ".dotfiles")
	dotfilesBinDir := filepath.Join(dotfilesDir, "bin")

	println("Creating symlink: ", cwd, " -> ", dotfilesDir)
	helpers.GenerateSymlink(cwd, dotfilesDir)

	println("Adding ", dotfilesBinDir, " to PATH")
	helpers.AddToEnvPath(helpers.ScopeMachine, dotfilesBinDir)

	helpers.PressAnyKeyOrWaitToExit()
}
