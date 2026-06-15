//go:build !windows

package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/logrusorgru/aurora/v4"
)

func inheritOwnership(target string, createdDirs []string) {
	for _, dir := range createdDirs {
		inheritFromParent(dir, false)
	}

	inheritFromParent(target, true)
}

func inheritFromParent(path string, isSymlink bool) {
	parent := filepath.Dir(path)

	info, err := os.Stat(parent)
	if err != nil {
		fmt.Println(aurora.Red("UNEXPECTED: Error reading parent directory: " + parent))
		return
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return
	}

	if chownErr := os.Lchown(path, int(stat.Uid), int(stat.Gid)); chownErr != nil {
		fmt.Println(aurora.Red("UNEXPECTED: Error setting owner: " + path))
	}

	if isSymlink {
		return
	}

	if chmodErr := os.Chmod(path, info.Mode().Perm()); chmodErr != nil {
		fmt.Println(aurora.Red("UNEXPECTED: Error setting permissions: " + path))
	}
}
