package main

import (
	helpers "dotfiles/src"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func isTaskEnabled(taskName string) bool {
	output, err := exec.Command("schtasks", "/Query", "/TN", "\\"+taskName, "/FO", "LIST").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "Status:") && !strings.Contains(string(output), "Disabled")
}

func enableTask(taskName string) {
	cmd := exec.Command("schtasks", "/Change", "/TN", "\\"+taskName, "/Enable")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return
	}
}

func disableTask(taskName string) {
	cmd := exec.Command("schtasks", "/Change", "/TN", "\\"+taskName, "/Disable")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

func runScript(scriptName string) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	script := filepath.Join(cwd, "./slack/"+scriptName)
	start := exec.Command("cscript", "//nologo", script)
	start.Stdout = os.Stdout
	start.Stderr = os.Stderr
	if err := start.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "cscript failed:", err)
	}
}

func renderSlackStatus() {
	if isTaskEnabled("task-slack-force-start") {
		println("Slack: \033[32mAlways On\033[0m")
	} else if isTaskEnabled("task-slack-start") {
		println("Slack: \033[33mWork Time\033[0m")
	} else {
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
	helpers.EnsureAdminExecution()

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
				slackAlwaysStart()
			} else if selected == 1 {
				slackStartOnTime()
			} else {
				slackDisable()
			}

			helpers.PressAnyKeyOrWaitToExit()
			return true, nil
		case keys.Escape, keys.CtrlC:
			clearScreen()
			renderSlackStatus()
			println("Exiting...")

			helpers.PressAnyKeyOrWaitToExit()
			return true, nil
		}

		render(selected, false)
		return false, nil
	})
}

func slackAlwaysStart() {
	disableTask("task-slack-start")
	disableTask("task-slack-exit")

	enableTask("task-slack-force-start")
	runScript("force-start.vbs")
}

func slackStartOnTime() {
	disableTask("task-slack-force-start")

	enableTask("task-slack-start")
	enableTask("task-slack-exit")
	runScript("start.vbs")
	runScript("exit.vbs")
}

func slackDisable() {
	disableTask("task-slack-force-start")
	disableTask("task-slack-start")
	enableTask("task-slack-exit")
	runScript("force-exit.vbs")
}
