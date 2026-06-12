//go:build !windows

package helpers

import (
	"os/exec"
	"syscall"
)

func applySysProcAttr(cmd *exec.Cmd, detached bool) {
	if detached {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
	}
}
