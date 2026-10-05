package helpers

import (
	"fmt"
	"runtime"
)

func PowerDown() error {
	return changePowerState(false)
}

func PowerRestart() error {
	return changePowerState(true)
}

func changePowerState(restart bool) error {
	switch runtime.GOOS {
	case "darwin":
		action := "shut down"
		script := `tell application "loginwindow" to «event aevtrsdn»`
		if restart {
			action = "restart"
			script = `tell application "loginwindow" to «event aevtrrst»`
		}

		if err := ExecNativeCommand([]string{"osascript", "-e", script}); err != nil {
			return fmt.Errorf("failed to %s macOS: %w", action, err)
		}

		return nil
	case "windows":
		action := "/s"
		if restart {
			action = "/r"
		}

		if err := ExecNativeCommand([]string{"shutdown.exe", action, "/t", "0"}); err != nil {
			return fmt.Errorf("failed to change Windows power state: %w", err)
		}

		return nil
	case "linux":
		action := "poweroff"
		if restart {
			action = "reboot"
		}

		if err := ExecNativeCommand([]string{"systemctl", action}, ExecCommandOptions{AsAdmin: true}); err != nil {
			return fmt.Errorf("failed to %s Linux: %w", action, err)
		}

		return nil
	default:
		return fmt.Errorf("power actions are not supported on %s", runtime.GOOS)
	}
}
