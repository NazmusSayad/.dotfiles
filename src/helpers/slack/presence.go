package slack

import (
	"fmt"
	"slices"
	"time"

	"dotfiles/src/constants"
	helpers "dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

type SlackStatus string

const (
	SlackStatusAlways   SlackStatus = "always"
	SlackStatusWorkTime SlackStatus = "work-hours"
	SlackStatusDisabled SlackStatus = "disabled"
)

func isWorkTime(config OutputSlackConfig) (bool, string) {
	now := time.Now().In(constants.RECOMMENDED_TIMEZONE)
	weekday := now.Weekday()
	hour := now.Hour()

	if slices.Contains(config.OfficeTimeWeekends, weekday) {
		return false, "Weekend (" + weekday.String() + ")"
	}

	todayHash := GenerateOffDaysHash(now.Month(), now.Day())
	if slices.Contains(config.OfficeTimeOffDays, todayHash) {
		return false, "Off day (" + todayHash + ")"
	}

	isWorkTime := hour >= config.OfficeTimeStart && hour < config.OfficeTimeFinish
	return isWorkTime, "Not Office Hours"
}

func SlackUpdatePresence(status SlackStatus) {
	for workspace, config := range ReadSlackConfig() {
		switch status {
		case SlackStatusAlways:
			fmt.Println("> " + workspace + ": Setting active...")
			setPresence(workspace, "active")

		case SlackStatusDisabled:
			fmt.Println("> " + workspace + ": Setting away...")
			setPresence(workspace, "away")

		case SlackStatusWorkTime:
			isWorkTime, reason := isWorkTime(config)

			if isWorkTime {
				fmt.Println("> " + workspace + ": Currently Office Hours, setting active...")
				setPresence(workspace, "active")
			} else {
				fmt.Println("> " + workspace + ": Currently " + reason + ", setting away...")
				setPresence(workspace, "away")
			}

		default:
			panic("invalid slack status: " + string(status))
		}
	}
}

func setPresence(workspace string, presence string) {
	err := helpers.ExecNativeCommand(
		[]string{"agent-slack", "user", "presence", presence, "--workspace", workspace},
		helpers.ExecCommandOptions{Silent: true},
	)
	if err != nil {
		fmt.Println(aurora.Red("Error: Failed to set " + workspace + " to " + presence))
	}
}
