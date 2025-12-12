package helpers

import (
	"os"
	"os/exec"
)

func ExecWithNativeOutput(exe string, args ...string) error {
	cmd := exec.Command(exe, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func ExecWithNativeOutputAndExit(exe string, args ...string) {
	if err := ExecWithNativeOutput(exe, args...); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			os.Exit(1)
		}
	}
}
