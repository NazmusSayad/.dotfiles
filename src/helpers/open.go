package helpers

import (
	"fmt"
	"os"
	"runtime"
)

func Open(target string, options ...ExecCommandOptions) error {
	if runtime.GOOS == "windows" {
		return ExecNativeCommand([]string{"rundll32", "url.dll,FileProtocolHandler", target}, options...)
	}

	if runtime.GOOS == "darwin" {
		return ExecNativeCommand([]string{"open", target}, options...)
	}

	if runtime.GOOS == "linux" {
		return ExecNativeCommand([]string{"xdg-open", target}, options...)
	}

	fmt.Fprintln(os.Stderr, "unsupported platform: "+runtime.GOOS)
	os.Exit(1)
	return nil
}
