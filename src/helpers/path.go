package helpers

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(input string) string {
	if strings.HasPrefix(input, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		input = filepath.Join(homeDir, input[1:])
	}

	if strings.HasPrefix(input, "@/") {
		dotfilesPath := os.Getenv("DOTFILES_DIR")
		if dotfilesPath == "" {
			panic("DOTFILES_DIR environment variable is not set")
		}

		input = filepath.Join(dotfilesPath, input[1:])
	}

	return os.ExpandEnv(input)
}
