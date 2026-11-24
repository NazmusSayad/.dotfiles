package helpers

import (
	"fmt"
	"os"
	"runtime"

	"github.com/logrusorgru/aurora/v4"
)

type ClaudeDevice struct {
	OS  string
	WSL bool
}

const CLAUDE_HOST_DEVICE = "windows"

var CLAUDE_DEVICES = map[string]ClaudeDevice{
	"windows": {OS: "windows"},
	"wsl":     {OS: "linux", WSL: true},
	"linux":   {OS: "linux"},
	"mac":     {OS: "darwin"},
}

func DetectClaudeDevice() string {
	for id, device := range CLAUDE_DEVICES {
		if device.OS == runtime.GOOS && device.WSL == IsWSL() {
			return id
		}
	}

	fmt.Println(aurora.Red("Error: Unknown device: " + runtime.GOOS))
	os.Exit(1)
	return ""
}
