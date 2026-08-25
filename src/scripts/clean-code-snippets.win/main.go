package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

var EDITORS_LIST = []Editor{
	{
		Name: "VSCode",
		Path: getCodeBasedRoot("code"),
	},
	{
		Name: "VSCode Insiders",
		Path: getCodeBasedRoot("code-insiders"),
	},
	{
		Name: "Cursor",
		Path: getCursorBasedRoot("cursor"),
	},
}

func getCodeBasedRoot(bin string) string {
	source, err := exec.LookPath(bin)
	if err != nil {
		return ""
	}

	rootDir := filepath.Join(source, "../../")
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return ""
	}

	sourceHash := ""
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "bin" {
			sourceHash = entry.Name()
			break
		}
	}
	if sourceHash == "" {
		return ""
	}

	return filepath.Join(rootDir, sourceHash)
}

func getCursorBasedRoot(bin string) string {
	binPath, err := exec.LookPath(bin)
	if err != nil {
		return ""
	}

	return filepath.Join(binPath, "../../../../")
}

func cleanupSnippets(editor Editor) {
	fmt.Println(">", aurora.Blue(editor.Name))

	if !utils.IsFileExists(editor.Path) {
		fmt.Println("Root dir not found:", aurora.Faint(editor.Path))
		return
	}

	extensionsDir := filepath.Clean(editor.Path + "/resources/app/extensions")
	if !utils.IsFileExists(extensionsDir) {
		fmt.Println("Extensions dir not found:", aurora.Faint(editor.Path))
		return
	}

	var files []string
	err := filepath.WalkDir(extensionsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".code-snippets") {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		fmt.Println(aurora.Red("Failed to walk through extensions dir"))
		return
	}

	for _, file := range files {
		err := os.WriteFile(file, []byte("{}"), 0o644)
		if err != nil {
			fmt.Println("! 	Failed clearing ", file, ": ", err)
			return
		}

		fmt.Println("✔️", aurora.Green(filepath.Base(file)))
	}
}

func main() {
	for _, editor := range EDITORS_LIST {
		cleanupSnippets(editor)
		fmt.Println()
	}
}

type Editor struct {
	Name string
	Path string
}
