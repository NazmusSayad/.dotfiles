package slack_helpers

import (
	"time"
)

var workTimeStart = time.Date(2025, 12, 5, 9, 0, 0, 0, time.Local)
var workTimeEnd = time.Date(2025, 12, 5, 17, 0, 0, 0, time.Local)

func SlackLaunch(status string) {
	if status == "always" {
		SlackApplicationStart()
	} else if status == "never" {
		SlackApplicationStop()
	} else if status == "work-hours" {
		currentTime := time.Now()

		println("Current time:", currentTime.Format("2006-01-02 15:04:05"))
		println("Work time start:", workTimeStart.Format("2006-01-02 15:04:05"))
		println("Work time end:", workTimeEnd.Format("2006-01-02 15:04:05"))

		if currentTime.After(workTimeStart) && currentTime.Before(workTimeEnd) {
			println("Starting Slack...")
			SlackApplicationStart()
		} else {
			println("Stopping Slack...")
			SlackApplicationStop()
		}
	}
}
