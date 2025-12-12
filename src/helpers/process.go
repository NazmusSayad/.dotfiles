package helpers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

func GetParentProcessName() (string, error) {
	ppid := os.Getppid()
	cmd := exec.Command(
		"cmd", "/C",
		"tasklist",
		"/FI", fmt.Sprintf("PID eq %d", ppid),
		"/FO", "CSV",
		"/NH",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	line := strings.TrimSpace(string(out))
	if line == "" || strings.HasPrefix(line, "INFO:") {
		return "", fmt.Errorf("no result")
	}

	r := csv.NewReader(strings.NewReader(line))
	rec, err := r.Read()
	if err != nil || len(rec) == 0 {
		return "", fmt.Errorf("parse error")
	}

	return strings.TrimSpace(rec[0]), nil
}

func IsStandaloneProcess() bool {
	parentProcessName, _ := GetParentProcessName()
	return parentProcessName == "explorer.exe"
}

func PressAnyKeyOrWaitToExit() {
	if !IsStandaloneProcess() {
		os.Exit(0)
	}

	const totalSeconds = 5
	fmt.Printf("Press any key to exit, or wait %d seconds...", totalSeconds)
	done := make(chan struct{}, 1)
	var h uintptr
	var orig uint32

	if runtime.GOOS == "windows" {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		getStdHandle := kernel32.NewProc("GetStdHandle")
		getConsoleMode := kernel32.NewProc("GetConsoleMode")
		setConsoleMode := kernel32.NewProc("SetConsoleMode")

		const STD_INPUT_HANDLE = uintptr(^uint32(10) + 1)
		const ENABLE_ECHO_INPUT = 0x0004
		const ENABLE_LINE_INPUT = 0x0002

		h, _, _ = getStdHandle.Call(STD_INPUT_HANDLE)

		var mode uint32
		_, _, _ = getConsoleMode.Call(h, uintptr(unsafe.Pointer(&mode)))
		orig = mode
		mode &^= (ENABLE_ECHO_INPUT | ENABLE_LINE_INPUT)
		_, _, _ = setConsoleMode.Call(h, uintptr(mode))

		go func() {
			b := make([]byte, 1)
			_, _ = os.Stdin.Read(b)
			done <- struct{}{}
		}()
	} else {
		go func() {
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadByte()
			done <- struct{}{}
		}()
	}

	deadline := time.Now().Add(totalSeconds * time.Second)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			if runtime.GOOS == "windows" {
				kernel32 := syscall.NewLazyDLL("kernel32.dll")
				setConsoleMode := kernel32.NewProc("SetConsoleMode")
				_, _, _ = setConsoleMode.Call(h, uintptr(orig))
			}
			println()
			os.Exit(0)
		case <-ticker.C:
			remaining := int(math.Ceil(time.Until(deadline).Seconds()))
			if remaining <= 0 {
				if runtime.GOOS == "windows" {
					kernel32 := syscall.NewLazyDLL("kernel32.dll")
					setConsoleMode := kernel32.NewProc("SetConsoleMode")
					_, _, _ = setConsoleMode.Call(h, uintptr(orig))
				}
				println()
				os.Exit(0)
			}
			fmt.Printf("\rPress any key to exit, or wait %d seconds...", remaining)
		}
	}
}
