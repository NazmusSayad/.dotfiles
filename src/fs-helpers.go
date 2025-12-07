package helpers

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/jsonc"
)

func ResolvePath(input string) string {
	if strings.HasPrefix(input, ".") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			println("Error: failed to get user home directory:", err)
			os.Exit(1)
		}

		dotfilesPath := filepath.Join(homeDir, ".dotfiles")

		if _, err := os.Stat(dotfilesPath); err == nil {
			input = filepath.Join(dotfilesPath, input)
		} else {
			println("Error: .dotfiles directory not found.")
			println("Please run __install-dotfiles.cmd to install the dotfiles.")
			os.Exit(1)
		}
	}

	return os.ExpandEnv(input)
}

func ReadDotfilesConfigJSONC(path string) ([]byte, error) {
	resolvedPath := ResolvePath(path)
	println("JSON:", resolvedPath)

	f, err := os.Open(resolvedPath)
	if err != nil {
		println("JSON: failed to open file")
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		println("JSON: failed to read file")
		return nil, err
	}

	return jsonc.ToJSON(data), nil
}
