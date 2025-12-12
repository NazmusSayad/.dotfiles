package slack_helpers

import (
	"time"
)

var bangladeshTZ = time.FixedZone("GMT+6", 6*60*60)

type SlackStatus string

const (
	SlackStatusAlways   SlackStatus = "always"
	SlackStatusWorkTime SlackStatus = "work-hours"
	SlackStatusDisabled SlackStatus = "disabled"
)

func isWorkTime() bool {
	now := time.Now().In(bangladeshTZ)
	weekday := now.Weekday()
	hour := now.Hour()

	if weekday == time.Friday || weekday == time.Saturday {
		return false
	}

	return hour >= 6 && hour < 20
}

func SlackLaunch(status SlackStatus) {
	switch status {
	case SlackStatusAlways:
		println("> Starting Slack...")
		SlackApplicationStart()

	case SlackStatusDisabled:
		println("> Stopping Slack...")
		SlackApplicationStop()

	case SlackStatusWorkTime:
		if isWorkTime() {
			println("> Currently in work time, starting Slack...")
			SlackApplicationStart()
		} else {
			println("> Currently not in work time, stopping Slack...")
			SlackApplicationStop()
		}
	}
}
