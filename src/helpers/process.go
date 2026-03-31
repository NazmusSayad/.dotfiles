package helpers

import (
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

func FindOtherPIDByExecutablePath(executablePath string, currentPID int32) (int32, error) {
	processes, err := process.Processes()
	if err != nil {
		return 0, err
	}

	for _, p := range processes {
		if p.Pid == currentPID {
			continue
		}

		processExecutablePath, err := p.Exe()
		if err != nil {
			continue
		}

		if isSameExecutablePath(processExecutablePath, executablePath) {
			return p.Pid, nil
		}
	}

	return 0, nil
}

func TerminateProcessTree(pid int32) error {
	if pid <= 0 {
		return nil
	}

	p, err := process.NewProcess(pid)
	if err != nil {
		return nil
	}

	children, err := p.Children()
	if err == nil {
		for _, child := range children {
			_ = TerminateProcessTree(child.Pid)
		}
	}

	err = p.Terminate()
	if err != nil {
		if err := p.Kill(); err != nil {
			return err
		}

		return nil
	}

	for range 20 {
		running, err := p.IsRunning()
		if err != nil || !running {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	if err := p.Kill(); err != nil {
		return err
	}

	return nil
}

func isSameExecutablePath(a string, b string) bool {
	resolvedA, err := filepath.Abs(a)
	if err != nil {
		resolvedA = a
	}

	resolvedB, err := filepath.Abs(b)
	if err != nil {
		resolvedB = b
	}

	if runtime.GOOS == "windows" {
		return strings.EqualFold(filepath.Clean(resolvedA), filepath.Clean(resolvedB))
	}

	return filepath.Clean(resolvedA) == filepath.Clean(resolvedB)
}
