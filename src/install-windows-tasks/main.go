package main

import (
	"dotfiles/src/constants"
	"dotfiles/src/helpers"
	slack_helpers "dotfiles/src/helpers/slack"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	helpers.EnsureAdminExecution()

	tasks := map[string]string{
		constants.STARTUP_TASK_SLACK.TaskName:   generateSlackTaskXML(resolveProgramPath(constants.STARTUP_TASK_SLACK.Program)),
		constants.STARTUP_TASK_WINDOWS.TaskName: generateWindowsTaskXML(resolveProgramPath(constants.STARTUP_TASK_WINDOWS.Program)),
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

func resolveProgramPath(program string) string {
	return helpers.ResolvePath("@/" + constants.BUILD_TASKS_RUNNER_DIR + "/" + program + ".vbs")
}

func generateSlackTaskXML(program string) string {
	config := slack_helpers.ReadSlackConfig()

	return `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.4" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>` + time.Now().In(constants.RECOMMENDED_TIMEZONE).Format(time.RFC3339) + `</Date>
    <Author>Nazmus Sayad</Author>
    <URI>\__SLACK-STARTUP-TASK___</URI>
  </RegistrationInfo>
  <Triggers>
    <BootTrigger>
      <Enabled>true</Enabled>
    </BootTrigger>
    <LogonTrigger>
      <Enabled>true</Enabled>
    </LogonTrigger>
    <EventTrigger>
      <Enabled>true</Enabled>
      <Subscription>&lt;QueryList&gt;&lt;Query Id="0" Path="System"&gt;&lt;Select
        Path="System"&gt;*[System[Provider[@Name='Microsoft-Windows-Power-Troubleshooter'] and
        EventID=1]]&lt;/Select&gt;&lt;/Query&gt;&lt;/QueryList&gt;</Subscription>
    </EventTrigger>
    ` + generateCalendarTriggerXML(config.OfficeTimeStart) + `
    ` + generateCalendarTriggerXML(config.OfficeTimeStart+1) + `
    ` + generateCalendarTriggerXML(config.OfficeTimeStart+2) + `
    ` + generateCalendarTriggerXML(config.OfficeTimeFinish) + `
    ` + generateCalendarTriggerXML(config.OfficeTimeFinish+1) + `
    ` + generateCalendarTriggerXML(config.OfficeTimeFinish+2) + `
  </Triggers>
  <Principals>
    <Principal id="Author">
      <GroupId>S-1-5-32-544</GroupId>
      <RunLevel>LeastPrivilege</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>true</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <DisallowStartOnRemoteAppSession>false</DisallowStartOnRemoteAppSession>
    <UseUnifiedSchedulingEngine>true</UseUnifiedSchedulingEngine>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>` + program + `</Command>
    </Exec>
  </Actions>
</Task>`
}

func generateWindowsTaskXML(program string) string {
	return `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.4" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>` + time.Now().Format(time.RFC3339) + `</Date>
    <Author>Nazmus Sayad</Author>
    <URI>\__OS STARTUP TASK__</URI>
  </RegistrationInfo>
  <Triggers>
    <BootTrigger>
      <Enabled>true</Enabled>
    </BootTrigger>
    <LogonTrigger>
      <Enabled>true</Enabled>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <GroupId>S-1-5-32-544</GroupId>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>true</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>true</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <DisallowStartOnRemoteAppSession>false</DisallowStartOnRemoteAppSession>
    <UseUnifiedSchedulingEngine>true</UseUnifiedSchedulingEngine>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <Priority>7</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>` + program + `</Command>
    </Exec>
  </Actions>
</Task>`
}

func generateCalendarTriggerXML(hour int) string {
	return `<CalendarTrigger>
<StartBoundary>` + formatTimeForCalendarTrigger(hour, 1) + `</StartBoundary>
<Enabled>true</Enabled>
<ScheduleByDay>
	<DaysInterval>1</DaysInterval>
</ScheduleByDay>
</CalendarTrigger>`
}

func formatTimeForCalendarTrigger(hour int, minute int) string {
	return time.Date(1970, 1, 1, hour, minute, 0, 0, constants.RECOMMENDED_TIMEZONE).Format(time.RFC3339)
}
