package helpers

import (
	"strings"
	"syscall"
	"unsafe"
)

var (
	shell32        = syscall.NewLazyDLL("shell32.dll")
	procShellExecW = shell32.NewProc("ShellExecuteW")
)

func DetachedElevate(exe string, args ...string) error {
	verb := syscall.StringToUTF16Ptr("runas")
	path := syscall.StringToUTF16Ptr(exe)

	var params *uint16
	if len(args) > 0 {
		p := strings.Join(args, " ")
		params = syscall.StringToUTF16Ptr(p)
	}

	// 0 = SW_HIDE
	r, _, _ := procShellExecW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(params)),
		0,
		0,
	)
	if r <= 32 {
		return syscall.Errno(r)
	}
	return nil
}
