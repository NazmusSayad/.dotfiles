package helpers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
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

func ReadDotfilesConfigJSONC(path string) ([]byte, error) {
	resolvedPath := ResolvePath(path)
	fmt.Println(aurora.Faint("JSON: " + resolvedPath))

	f, err := os.Open(resolvedPath)
	if err != nil {
		fmt.Println(aurora.Red("JSON: failed to open file"))
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(aurora.Red("JSON: failed to read file"))
		return nil, err
	}

	return jsonc.ToJSON(data), nil
}

func IsFileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	_ = fi
	return true
}
