package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	helpers "dotfiles/src"
	slack_helpers "dotfiles/src/slack/helpers"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func getSlackStatusFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".slack-status")
}

func readSlackStatus() string {
	data, err := os.ReadFile(getSlackStatusFilePath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func writeSlackStatus(status string) {
	os.WriteFile(getSlackStatusFilePath(), []byte(status), 0644)
}

func renderSlackStatus() {
	switch readSlackStatus() {
	case "always":
		println("Slack: \033[32mAlways On\033[0m")
	case "work-hours":
		println("Slack: \033[33mWork Time\033[0m")
	default:
		println("Slack: \033[31mDisabled\033[0m")
	}
}

var options = []string{"Enable Always", "Enable on Work Time", "Disable"}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func render(selected int, final bool) {
	clearScreen()
	renderSlackStatus()
	println()

	for i, opt := range options {
		if i == selected {
			println("  \033[36m❯ " + opt + "\033[0m")
		} else {
			println("    " + opt)
		}
	}
	println()

	if final {
		println("Updating slack status to: \033[36m" + options[selected] + "\033[0m")
	} else {
		println("\033[90mUse ↑/↓ to navigate, Enter to select, Esc to exit\033[0m")
	}
}

func main() {
	selected := 0
	render(selected, false)

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Up:
			if selected > 0 {
				selected--
			}

		case keys.Down:
			if selected < len(options)-1 {
				selected++
			}

		case keys.Enter:
			render(selected, true)
			println()

			if selected == 0 {
				writeSlackStatus("always")
				slack_helpers.SlackLaunch("always")
				helpers.PressAnyKeyOrWaitToExit()
			} else if selected == 1 {
				writeSlackStatus("work-hours")
				slack_helpers.SlackLaunch("work-hours")
				helpers.PressAnyKeyOrWaitToExit()
			} else {
				writeSlackStatus("never")
				slack_helpers.SlackLaunch("never")
				helpers.PressAnyKeyOrWaitToExit()
			}

		case keys.Escape, keys.CtrlC:
			clearScreen()
			renderSlackStatus()
			println("Exiting...")
			helpers.PressAnyKeyOrWaitToExit()
		}

		render(selected, false)
		return false, nil
	})
}
