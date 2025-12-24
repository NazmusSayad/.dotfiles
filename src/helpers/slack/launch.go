package slack_helpers

import (
	"dotfiles/src/constants"
	"fmt"
	"slices"
	"time"
)

type SlackStatus string

const (
	SlackStatusAlways   SlackStatus = "always"
	SlackStatusWorkTime SlackStatus = "work-hours"
	SlackStatusDisabled SlackStatus = "disabled"
)

func isWorkTime() (bool, string) {
	config := ReadSlackConfig()

	now := time.Now().In(constants.RECOMMENDED_TIMEZONE)
	weekday := now.Weekday()
	hour := now.Hour()

	if slices.Contains(config.OfficeTimeWeekend, weekday) {
		return false, "Weekend (" + weekday.String() + ")"
	}

	todayHash := GenerateOffDaysHash(now.Month(), now.Day())
	if slices.Contains(config.OfficeTimeOffDays, todayHash) {
		return false, "Off day (" + todayHash + ")"
	}

	isWorkTime := hour >= config.OfficeTimeStart && hour < config.OfficeTimeFinish
	return isWorkTime, "Not Office Hours"
}

func SlackLaunch(status SlackStatus) {
	switch status {
	case SlackStatusAlways:
		fmt.Println("> Starting Slack...")
		SlackApplicationStart()

	case SlackStatusDisabled:
		fmt.Println("> Stopping Slack...")
		SlackApplicationStop()

	case SlackStatusWorkTime:
		isWorkTime, reason := isWorkTime()

		if isWorkTime {
			fmt.Println("> Currently Office Hours, starting Slack...")
			SlackApplicationStart()
		} else {
			fmt.Println("> Currently " + reason + ", stopping Slack...")
			SlackApplicationStop()
		}
	}
}
