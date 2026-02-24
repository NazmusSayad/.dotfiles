package utils

import "strconv"

func HourToAmPm(h int) string {
	switch {
	case h == 0:
		return "12AM"
	case h == 12:
		return "12PM"
	case h < 12:
		return strconv.Itoa(h) + "AM"
	default:
		return strconv.Itoa(h-12) + "PM"
	}
}
