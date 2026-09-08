//go:build !windows

package helpers

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/logrusorgru/aurora/v4"
)

func ApplyUserOwnership(paths ...string) {
	if os.Geteuid() != 0 {
		return
	}

	uid, gid, err := invokingUser()
	if err != nil {
		fmt.Println(aurora.Red("UNEXPECTED: " + err.Error()))
		return
	}

	for _, path := range paths {
		walkErr := filepath.WalkDir(path, func(current string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			return os.Lchown(current, uid, gid)
		})

		if walkErr != nil && !os.IsNotExist(walkErr) {
			fmt.Println(aurora.Red("UNEXPECTED: Error setting owner: " + path))
		}
	}
}

func ValidateInvokingUser() error {
	if os.Geteuid() != 0 {
		return nil
	}

	_, _, err := invokingUser()
	return err
}

func invokingUser() (int, int, error) {
	uid, uidErr := strconv.Atoi(os.Getenv("SUDO_UID"))
	gid, gidErr := strconv.Atoi(os.Getenv("SUDO_GID"))

	if uidErr != nil || gidErr != nil {
		return 0, 0, fmt.Errorf("cannot determine invoking user: SUDO_UID and SUDO_GID must be set")
	}
	if uid <= 0 || gid < 0 {
		return 0, 0, fmt.Errorf("cannot use invoking user UID %d and GID %d", uid, gid)
	}

	return uid, gid, nil
}

func inheritDirOwnership(dirs []string) {
	for _, dir := range dirs {
		inheritFromParent(dir, false)
	}
}

func inheritOwnership(target string, createdDirs []string) {
	inheritDirOwnership(createdDirs)
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
