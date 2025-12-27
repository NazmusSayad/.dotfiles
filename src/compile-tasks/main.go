package main

import (
	"dotfiles/src/constants"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

var WINDOWS_TASKS = []constants.WindowsTask{
	constants.STARTUP_TASK_SLACK,
	constants.STARTUP_TASK_WINDOWS,
}

func main() {
	if err := os.MkdirAll(constants.BUILD_TASKS_RUNNER_DIR, 0755); err != nil {
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

		filename := filepath.Join(constants.BUILD_TASKS_RUNNER_DIR, runner.Program+".vbs")

		if err := os.WriteFile(filename, []byte(script), 0644); err != nil {
			fmt.Println(aurora.Red("Error: failed to write " + filename + ": " + err.Error()))
			os.Exit(1)
		}

		fmt.Println(aurora.Faint("> " + filename))
	}
}

func generateVbsScriptAsUser(program string) string {
	parts := []string{
		`CreateObject("WScript.Shell").Run`,
		`"` + program + `",`, // program
		`0,`,                 // hide window
		`True`,               // wait for completion
	}

	return strings.Join(parts, " ")
}

func generateVbsScriptAsAdmin(program string) string {
	parts := []string{
		`CreateObject("Shell.Application").ShellExecute`,
		`"` + program + `",`, // program
		`"",`,                // arguments
		`"",`,                // working directory
		`"runas",`,           // verb
		`0`,                  // hide window
	}
	return strings.Join(parts, " ")
}
