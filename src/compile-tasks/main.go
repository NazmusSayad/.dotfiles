package main

import (
	"dotfiles/src/constants"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

const USER_BIN_DIR = "$USERPROFILE/" + constants.SCRIPTS_BUILD_BIN_DIR

var WINDOWS_TASKS = []constants.WindowsTask{
	constants.STARTUP_TASK_SLACK,
	constants.STARTUP_TASK_WINDOWS,
}

func main() {
	if err := os.MkdirAll(constants.TASKS_RUNNER_BUILD_DIR, 0755); err != nil {
		fmt.Println(aurora.Red("Error: failed to create tasks directory: " + err.Error()))
		os.Exit(1)
	}

	for _, runner := range WINDOWS_TASKS {
		script := ""
		program := runner.Program + ".exe"

		if runner.Mode == constants.WindowsTaskModeAdmin {
			script = generateVbsScriptAsAdmin(program)
		} else {
			script = generateVbsScriptAsUser(program)
		}

		filename := filepath.Join(constants.TASKS_RUNNER_BUILD_DIR, runner.Program+".vbs")

		if err := os.WriteFile(filename, []byte(script), 0644); err != nil {
			fmt.Println(aurora.Red("Error: failed to write " + filename + ": " + err.Error()))
			os.Exit(1)
		}

		fmt.Println(aurora.Faint("> " + filename))
	}
}

func generateVbsScriptAsUser(program string) string {
	lines := []string{
		"Set WshShell = CreateObject(\"WScript.Shell\")",
		"WshShell.Run \"\"\"" + program + "\"\"\", 0, False",
		"Set WshShell = Nothing",
	}
	return strings.Join(lines, "\n")
}

func generateVbsScriptAsAdmin(program string) string {
	lines := []string{
		"Set UAC = CreateObject(\"Shell.Application\")",
		"UAC.ShellExecute \"" + program + "\", \"\", \"\", \"runas\", 0",
		"Set UAC = Nothing",
	}
	return strings.Join(lines, "\n")
}
