package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func ResolvePath(input string) string {
	if strings.HasPrefix(input, "@/") {
		dotfilesPath := os.Getenv("DOTFILES_DIR")

		if dotfilesPath == "" {
			fmt.Println(aurora.Red(".dotfiles environment variable is not set."))
			os.Exit(1)
		}

		input = filepath.Join(dotfilesPath, input[1:])
	}

	return os.ExpandEnv(input)
}
