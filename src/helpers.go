package helpers

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/tidwall/jsonc"
)

type Scope string

const (
	ScopeUser    Scope = "User"
	ScopeMachine Scope = "Machine"
)

func execPsCommand(command string) (string, error) {
	cmd := exec.Command("powershell", "-c", command)
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("powershell command failed: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func ReadEnv(scope Scope, name string) (string, error) {
	return execPsCommand(
		fmt.Sprintf(`[System.Environment]::GetEnvironmentVariable("%s", [System.EnvironmentVariableTarget]::%s)`, name, scope),
	)
}

func WriteEnv(scope Scope, name, value string) (string, error) {
	return execPsCommand(
		fmt.Sprintf(`[System.Environment]::SetEnvironmentVariable("%s", "%s", [System.EnvironmentVariableTarget]::%s)`, name, value, scope),
	)
}

func AddToEnvPath(scope Scope, paths ...string) (string, error) {
	existingPath, err := ReadEnv(scope, "PATH")
	if err != nil {
		return "", err
	}

	existingPathArray := strings.Split(existingPath, ";")
	var filteredPaths []string
	for _, p := range existingPathArray {
		if p != "" {
			filteredPaths = append(filteredPaths, p)
		}
	}

	pathSet := make(map[string]bool)
	var uniquePaths []string

	for _, p := range filteredPaths {
		if !pathSet[p] {
			pathSet[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	for _, p := range paths {
		if !pathSet[p] {
			pathSet[p] = true
			uniquePaths = append(uniquePaths, p)
		}
	}

	newPath := strings.Join(uniquePaths, ";")
	return WriteEnv(scope, "PATH", newPath)
}

func EnsureAdminExecution() {
	psCmd := `if (-not([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) { exit 1 }`
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psCmd)
	if err := cmd.Run(); err == nil {
		fmt.Println("Admin execution not required.")
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
			time.Sleep(2000)
			os.Exit(1)
		}

		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "Start-Process -FilePath '"+exePath+"' -Verb RunAs")
		if err := cmd.Run(); err != nil {
			fmt.Println("Failed to relaunch with elevated privileges.")
			fmt.Println("Press Enter to exit...")
			time.Sleep(2000)
			os.Exit(1)
		}

		fmt.Println("Relaunched with elevated privileges.")
		os.Exit(0)
	}
}

func ResolvePath(input string) string {
	if strings.HasPrefix(input, ".") {
		dotfilesPath := os.ExpandEnv("$HOME/.dotfiles")

		if _, err := os.Stat(dotfilesPath); err == nil {
			input = filepath.Join(dotfilesPath, input)
		} else {
			fmt.Println("Error: .dotfiles directory not found.")
			fmt.Println("Please run __install-dotfiles.cmd to install the dotfiles.")
			time.Sleep(2000)
			os.Exit(1)
		}
	}

	return os.ExpandEnv(input)
}

func PressAnyKeyOrWaitToExit() {
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
			fmt.Println()
			os.Exit(0)
		case <-ticker.C:
			remaining := int(math.Ceil(time.Until(deadline).Seconds()))
			if remaining <= 0 {
				if runtime.GOOS == "windows" {
					kernel32 := syscall.NewLazyDLL("kernel32.dll")
					setConsoleMode := kernel32.NewProc("SetConsoleMode")
					_, _, _ = setConsoleMode.Call(h, uintptr(orig))
				}
				fmt.Println()
				os.Exit(0)
			}
			fmt.Printf("\rPress any key to exit, or wait %d seconds...", remaining)
		}
	}
}

func ReadJsoncAsJson(path string) ([]byte, error) {
	fmt.Println("JSON:", path)

	f, err := os.Open(path)
	if err != nil {
		fmt.Println("JSON: failed to open file")
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("JSON: failed to read file")
		return nil, err
	}

	return jsonc.ToJSON(data), nil
}
