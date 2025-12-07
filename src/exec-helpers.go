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

func execPsCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("powershell command failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func EnsureAdminExecution() {
	psCmd := `if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	if err := cmd.Run(); err == nil {
		println("Admin execution not required.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintln(os.Stderr, "This program requires administrator privileges.")
		fmt.Fprintln(os.Stderr, "Trying to relaunch with elevated privileges...")

		exePath, e := os.Executable()
		if e != nil {
			_, _ = reader.ReadString('\n')

			println("Failed to get executable path.")
			os.Exit(1)
		}

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+exePath+"' -Verb RunAs")
		if err := cmd.Run(); err != nil {
			println("Failed to relaunch with elevated privileges.")
			println("Press Enter to exit...")
			os.Exit(1)
		}

		println("Relaunched with elevated privileges.")
		os.Exit(0)
	}
}

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

	return rec[0], nil
}

func IsEmbeddedInTerminal() bool {
	parentProcessName, err := GetParentProcessName()
	if err != nil || parentProcessName == "" || parentProcessName == "explorer.exe" {
		return false
	}

	return true
}

func PressAnyKeyOrWaitToExit() {
	if IsEmbeddedInTerminal() {
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
