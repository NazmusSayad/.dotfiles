package main

import (
	"dotfiles/src/helpers"
	slack "dotfiles/src/helpers/slack"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	helpers.EnsureAdminExecution()

	slackConfig := slack.ReadSlackConfig()
	runHidden, err := exec.LookPath("run-hidden.exe")
	if err != nil {
		panic(err)
	}

	tasks := map[string]string{
		"__SLACK STARTUP TASK___": helpers.GenerateWindowsTaskXML(helpers.WindowsTaskOptions{
			Author: "Nazmus Sayad",
			Mode:   helpers.WindowsTaskRunLevelLeastPrivilege,
			Triggers: []helpers.WindowsTaskTrigger{
				{Type: helpers.WindowsTaskTriggerTypeBoot},
				{Type: helpers.WindowsTaskTriggerTypeLogon},
				{Type: helpers.WindowsTaskTriggerTypeDaily, Hour: slackConfig.OfficeTimeStart, Minute: 0},
				{Type: helpers.WindowsTaskTriggerTypeDaily, Hour: slackConfig.OfficeTimeStart + 1, Minute: 0},
				{Type: helpers.WindowsTaskTriggerTypeDaily, Hour: slackConfig.OfficeTimeStart + 2, Minute: 0},
				{Type: helpers.WindowsTaskTriggerTypeDaily, Hour: slackConfig.OfficeTimeFinish, Minute: 0},
				{Type: helpers.WindowsTaskTriggerTypeDaily, Hour: slackConfig.OfficeTimeFinish + 1, Minute: 0},
				{Type: helpers.WindowsTaskTriggerTypeDaily, Hour: slackConfig.OfficeTimeFinish + 2, Minute: 0},
				{Type: helpers.WindowsTaskTriggerTypeEvent, EventString: `&lt;QueryList&gt;&lt;Query Id="0" Path="System"&gt;&lt;Select
        Path="System"&gt;*[System[Provider[@Name='Microsoft-Windows-Power-Troubleshooter'] and
        EventID=1]]&lt;/Select&gt;&lt;/Query&gt;&lt;/QueryList&gt;`},
			},
			Actions: []helpers.WindowsTaskAction{
				{
					Command:   runHidden,
					Arguments: []string{"slack-startup.exe"},
				},
			},
		}),

		"__WINDOWS STARTUP TASK__": helpers.GenerateWindowsTaskXML(helpers.WindowsTaskOptions{
			Author: "Nazmus Sayad",
			Mode:   helpers.WindowsTaskRunLevelHighestAvailable,
			Triggers: []helpers.WindowsTaskTrigger{
				{Type: helpers.WindowsTaskTriggerTypeBoot},
				{Type: helpers.WindowsTaskTriggerTypeLogon},
			},
			Actions: []helpers.WindowsTaskAction{
				{
					Command:   runHidden,
					Arguments: []string{"windows-startup.exe"},
				},
			},
		}),
	}

	for name, task := range tasks {
		tmp := os.TempDir()
		xmlPath := filepath.Join(tmp, "__"+name+".xml")

		if err := utils.WriteUTF16LE(xmlPath, task); err != nil {
			panic(err)
		}

		err := helpers.ExecNativeCommand(
			[]string{
				"schtasks",
				"/create",
				"/tn", name,
				"/xml", xmlPath,
				"/f",
			},
			helpers.ExecCommandOptions{
				Silent: true,
			},
		)

		_ = os.Remove(xmlPath)

		if err != nil {
			fmt.Println(aurora.Red("Failed to create task: ").String() + name)
		} else {
			fmt.Println(aurora.Green("Successfully created task: ").String() + name)
		}
	}
}
