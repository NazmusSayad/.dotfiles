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

func isWorkTime() bool {
	now := time.Now().In(constants.SLACK_TIMEZONE)
	weekday := now.Weekday()
	hour := now.Hour()

	if slices.Contains(constants.SLACK_OFFICE_HOUR_WEEKEND, weekday) {
		return false
	}

	return hour >= constants.SLACK_OFFICE_HOUR_START && hour < constants.SLACK_OFFICE_HOUR_FINISH
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
		if isWorkTime() {
			fmt.Println("> Currently in work time, starting Slack...")
			SlackApplicationStart()
		} else {
			fmt.Println("> Currently not in work time, stopping Slack...")
			SlackApplicationStop()
		}
	}
}
