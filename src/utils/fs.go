package utils

import "os"

func IsFileExists(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		return false
	}

	_ = fi
	return true
}
