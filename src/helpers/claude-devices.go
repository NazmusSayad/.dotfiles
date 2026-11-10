package helpers

import (
	"fmt"
	"os"
	"runtime"

	"github.com/logrusorgru/aurora/v4"
)

type ClaudeDevice struct {
	OS          string
	SSHEnv      string
	StoragePath string
}

const CLAUDE_HOST_DEVICE = "windows"

var CLAUDE_DEVICES = map[string]ClaudeDevice{
	"windows": {OS: "windows", SSHEnv: "SSH_WIN_ADDRESS", StoragePath: "F:/agents-store"},
	"wsl":     {OS: "linux", StoragePath: "/mnt/f/agents-store"},
	"mac":     {OS: "darwin"},
}

func DetectClaudeDevice() string {
	for id, device := range CLAUDE_DEVICES {
		if device.OS == runtime.GOOS {
			return id
		}
	}

	fmt.Println(aurora.Red("Error: Unknown device: " + runtime.GOOS))
	os.Exit(1)
	return ""
}

func ClaudeStorage(deviceID string) string {
	device := CLAUDE_DEVICES[deviceID]
	if device.StoragePath != "" {
		return device.StoragePath
	}

	hostDevice := CLAUDE_DEVICES[CLAUDE_HOST_DEVICE]
	address := os.Getenv(hostDevice.SSHEnv)
	if address == "" {
		fmt.Println(aurora.Red("Error: " + hostDevice.SSHEnv + " must be set"))
		os.Exit(1)
	}

	return "sftp:" + address + ":/" + hostDevice.StoragePath
}
