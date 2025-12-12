package helpers

import (
	"os/exec"
	"syscall"
)

func DetachedExec(exe string, args ...string) error {
	cmd := exec.Command(exe, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008}
	cmd.Start()
	return nil
}
