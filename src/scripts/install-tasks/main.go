package main

import (
	"dotfiles/src/constants"
	"dotfiles/src/helpers"
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
		"__SLACK TASK___":     generateSlackTaskXML(),
		"__OS STARTUP TASK__": generateWindowsTaskXML(),
	}

	for name, task := range tasks {
		tmp := os.TempDir()
		xmlPath := filepath.Join(tmp, "__"+name+".xml")

		if err := utils.WriteUTF16LE(xmlPath, task); err != nil {
			panic(err)
		}

		err := helpers.ExecNativeCommand(
			helpers.ExecCommandOptions{
				Command: "schtasks",
				Args: []string{
					"/create",
					"/tn", name,
					"/xml", xmlPath,
					"/f",
				},
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

func generateSlackTaskXML() string {
	actionCommand := "proxy-vbs.exe"
	actionArguments := "user slack-startup.exe"

	return `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.4" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>` + time.Now().In(constants.SLACK_TIMEZONE).Format(time.RFC3339) + `</Date>
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
    ` + generateCalendarTriggerXML(constants.SLACK_OFFICE_HOUR_START) + `
    ` + generateCalendarTriggerXML(constants.SLACK_OFFICE_HOUR_START+2) + `
    ` + generateCalendarTriggerXML(constants.SLACK_OFFICE_HOUR_FINISH) + `
    ` + generateCalendarTriggerXML(constants.SLACK_OFFICE_HOUR_FINISH+2) + `
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
      <Command>` + actionCommand + `</Command>
      <Arguments>` + actionArguments + `</Arguments>
    </Exec>
  </Actions>
</Task>`
}

func generateWindowsTaskXML() string {
	actionCommand := "proxy-vbs.exe"
	actionArguments := "admin windows-startup.exe"

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
      <Command>` + actionCommand + `</Command>
      <Arguments>` + actionArguments + `</Arguments>
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
	return time.Date(1970, 1, 1, hour, minute, 0, 0, constants.SLACK_TIMEZONE).Format(time.RFC3339)
}
