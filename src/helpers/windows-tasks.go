package helpers

import (
	"dotfiles/src/constants"
	"strings"
	"time"
)

type WindowsTaskRunLevel string

const (
	WindowsTaskRunLevelLeastPrivilege   WindowsTaskRunLevel = "LeastPrivilege"
	WindowsTaskRunLevelHighestAvailable WindowsTaskRunLevel = "HighestAvailable"
)

type WindowsTaskTriggerType string

const (
	WindowsTaskTriggerTypeDaily WindowsTaskTriggerType = "daily"
	WindowsTaskTriggerTypeEvent WindowsTaskTriggerType = "event"
	WindowsTaskTriggerTypeLogon WindowsTaskTriggerType = "logon"
	WindowsTaskTriggerTypeBoot  WindowsTaskTriggerType = "boot"
)

type WindowsTaskTrigger struct {
	Type WindowsTaskTriggerType
	Time time.Time

	Hour        int
	Minute      int
	EventString string
}

type WindowsTaskAction struct {
	Command   string
	Arguments []string
}

type WindowsTaskOptions struct {
	Name     string
	Author   string
	Mode     WindowsTaskRunLevel
	Triggers []WindowsTaskTrigger
	Actions  []WindowsTaskAction
}

func GenerateWindowsTaskXML(options WindowsTaskOptions) string {
	triggers := []string{}
	actions := []string{}

	for _, trigger := range options.Triggers {
		if trigger.Type == WindowsTaskTriggerTypeBoot {
			triggers = append(triggers, `<BootTrigger>
        <Enabled>true</Enabled>
      </BootTrigger>`,
			)
		}

		if trigger.Type == WindowsTaskTriggerTypeLogon {
			triggers = append(triggers, `<LogonTrigger>
        <Enabled>true</Enabled>
      </LogonTrigger>`,
			)
		}

		if trigger.Type == WindowsTaskTriggerTypeEvent {
			triggers = append(triggers, `<EventTrigger>
        <Enabled>true</Enabled>
        <Subscription>`+trigger.EventString+`</Subscription>
      </EventTrigger>`,
			)
		}

		if trigger.Type == "daily" {
			triggers = append(triggers, `<CalendarTrigger>
      <StartBoundary>`+formatTimeForCalendarTrigger(trigger.Hour, trigger.Minute)+`</StartBoundary>
      <Enabled>true</Enabled>
      <ScheduleByDay>
    	  <DaysInterval>1</DaysInterval>
      </ScheduleByDay>
    </CalendarTrigger>`,
			)
		}

	}

	for _, action := range options.Actions {

		argumentsString := ""
		if len(action.Arguments) > 0 {
			argumentsString = "\n      <Arguments>" + strings.Join(action.Arguments, " ") + "</Arguments>"
		}

		actions = append(actions, `<Exec>
      <Command>`+action.Command+`</Command>`+argumentsString+`
    </Exec>`,
		)
	}

	return `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.4" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Date>` + time.Now().In(constants.RECOMMENDED_TIMEZONE).Format(time.RFC3339) + `</Date>
    <Author>Nazmus Sayad</Author>
    <URI>\__SLACK-STARTUP-TASK___</URI>
  </RegistrationInfo>
  <Triggers>
    ` + strings.Join(triggers, "\n") + `
  </Triggers>
  <Principals>
    <Principal id="Author">
      <GroupId>S-1-5-32-544</GroupId>
      <RunLevel>` + string(options.Mode) + `</RunLevel>
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
    ` + strings.Join(actions, "\n") + `
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
