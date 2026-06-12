//go:build windows

package helpers

import (
	"os/exec"

	"golang.org/x/sys/windows"
)

func applySysProcAttr(cmd *exec.Cmd, detached bool) {
	if detached {
		cmd.SysProcAttr = &windows.SysProcAttr{
			HideWindow:       true,
			NoInheritHandles: true,
			CreationFlags:    windows.DETACHED_PROCESS | windows.CREATE_NEW_PROCESS_GROUP,
		}
	} else {
		cmd.SysProcAttr = &windows.SysProcAttr{
			HideWindow: true,
		}
	}
}
