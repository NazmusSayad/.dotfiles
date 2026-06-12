package helpers

import (
	"os"
	"path/filepath"
	"strings"

	"dotfiles/src/utils"
)

func ResolvePath(input string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(input, "~") {
		input = filepath.Join(homeDir, input[1:])
	}

	if strings.HasPrefix(input, "@/") {
		dotfilesPath := os.Getenv("DOTFILES_DIR")

		if dotfilesPath == "" {
			homeDirPath := filepath.Join(homeDir, ".dotfiles")

			if utils.IsFileExists(homeDirPath) {
				dotfilesPath = homeDirPath
			} else {
				panic("DOTFILES_DIR environment variable is not set")
			}
		}

		input = filepath.Join(dotfilesPath, input[1:])
	}

	return os.ExpandEnv(input)
}
