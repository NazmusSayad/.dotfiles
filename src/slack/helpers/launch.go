package slack_helpers

import (
	"time"
)

var bangladeshTZ = time.FixedZone("GMT+6", 6*60*60)

func isWorkTime() bool {
	now := time.Now().In(bangladeshTZ)
	weekday := now.Weekday()
	hour := now.Hour()

	if weekday == time.Friday || weekday == time.Saturday {
		return false
	}

	return hour >= 6 && hour < 20
}

func SlackLaunch(status string) {
	switch status {
	case "always":
		SlackApplicationStart()

	case "never":
		SlackApplicationStop()

	case "work-hours":
		now := time.Now().In(bangladeshTZ)
		println("Current time (BD):", now.Format("2006-01-02 15:04:05 Mon"))
		println("Work hours: 6AM-8PM, Off: Fri/Sat")

		if isWorkTime() {
			println("Starting Slack...")
			SlackApplicationStart()
		} else {
			println("Stopping Slack...")
			SlackApplicationStop()
		}
	}
}
