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
		if err := disableMacOSWindowRestore(); err != nil {
			return err
		}

		action := "shut down"
		if restart {
			action = "restart"
		}

		if err := ExecNativeCommand([]string{"osascript", "-e", `tell application "System Events" to ` + action}); err != nil {
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

func disableMacOSWindowRestore() error {
	keys := []string{"TALLogoutSavesState", "LoginwindowLaunchesRelaunchApps"}
	for _, key := range keys {
		if err := ExecNativeCommand([]string{"defaults", "write", "com.apple.loginwindow", key, "-bool", "false"}); err != nil {
			return fmt.Errorf("failed to disable macOS window restore: %w", err)
		}
	}

	return nil
}
