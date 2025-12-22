package constants

import "time"

const SLACK_OFFICE_HOUR_START = 6
const SLACK_OFFICE_HOUR_FINISH = 21

var SLACK_TIMEZONE = time.FixedZone("GMT+6", 6*60*60)

var SLACK_OFFICE_HOUR_WEEKEND = []time.Weekday{
	time.Friday,
	time.Saturday,
}
