package main

import (
	helpers "dotfiles/src/helpers"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Editor struct {
	Name           string
	Alias          string
	Path           string
	ExtensionsPath string
}

func findBinPath(bin string) string {
	cmd := "which"
	if runtime.GOOS == "windows" {
		cmd = "where"
	}
	out, err := exec.Command(cmd, bin).Output()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func getFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".code-snippets") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {
	editors := []Editor{
		{
			Name:           "VSCode",
			Alias:          "code",
			ExtensionsPath: "../../resources/app/extensions",
		},
		{
			Name:           "VSCode Insiders",
			Alias:          "code-insiders",
			ExtensionsPath: "../../resources/app/extensions",
		},
		{
			Name:           "Cursor",
			Alias:          "cursor",
			ExtensionsPath: "../../../../resources/app/extensions",
		},
	}

	for _, editor := range editors {
		binPath := findBinPath(editor.Alias)
		if binPath == "" {
			continue
		}

		binDir := filepath.Dir(binPath)
		extPath := filepath.Clean(binDir + editor.ExtensionsPath)

		if !helpers.IsFileExists(extPath) {
			println("Extensions path not found for ", editor.Name)
			println("Path: ", binDir)
			println("Ext: ", editor.ExtensionsPath)
			println("Resolved: ", extPath)
			println("")
			continue
		}

		files, err := getFiles(extPath)
		if err != nil {
			println("Failed reading files: ", err)
			continue
		}

		for _, file := range files {
			err := os.WriteFile(file, []byte("{}"), 0644)
			if err != nil {
				println("Failed clearing ", file, ": ", err)
				continue
			}
			println(editor.Name, " Cleared: ", filepath.Base(file))
		}
	}
}
