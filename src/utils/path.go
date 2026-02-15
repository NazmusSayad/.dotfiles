package utils

import "os/exec"

func IsCommandInPath(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
