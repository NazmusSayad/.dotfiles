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
	config := ReadSlackConfig()

	now := time.Now().In(constants.RECOMMENDED_TIMEZONE)
	weekday := now.Weekday()
	hour := now.Hour()

	if slices.Contains(config.OfficeTimeWeekend, weekday) {
		return false
	}

	return hour >= config.OfficeTimeStart && hour < config.OfficeTimeFinish
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
