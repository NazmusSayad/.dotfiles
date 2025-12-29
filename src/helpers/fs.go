package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func ResolvePath(input string) string {
	if strings.HasPrefix(input, "@") {
		dotfilesPath := os.Getenv("DOTFILES_DIR")

		if _, err := os.Stat(dotfilesPath); err == nil {
			input = filepath.Join(dotfilesPath, input[1:])
		} else {
			fmt.Println(aurora.Red("Error: .dotfiles directory not found."))
			fmt.Println(aurora.Yellow("Please run __install-dotfiles.cmd to install the dotfiles."))
			os.Exit(1)
		}
	}

	return os.ExpandEnv(input)
}

func IsFileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	_ = fi
	return true
}
