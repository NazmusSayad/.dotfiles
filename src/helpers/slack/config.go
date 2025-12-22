package slack_helpers

import (
	helpers "dotfiles/src/helpers"
	"encoding/json"
	"strings"
	"time"
)

type InputSlackConfig struct {
	OfficeTimeStart   int
	OfficeTimeFinish  int
	OfficeTimeWeekend []string
}

type OutputSlackConfig struct {
	OfficeTimeStart   int
	OfficeTimeFinish  int
	OfficeTimeWeekend []time.Weekday
}

func ReadSlackConfig() OutputSlackConfig {
	jsonBytes, err := helpers.ReadDotfilesConfigJSONC("./config/slack-status.jsonc")
	if err != nil {
		panic("Error reading slack config")
	}

	var input InputSlackConfig
	if err := json.Unmarshal(jsonBytes, &input); err != nil {
		panic("Error parsing slack config")
	}

	weekends := make([]time.Weekday, 0, len(input.OfficeTimeWeekend))
	for _, day := range input.OfficeTimeWeekend {
		switch strings.ToLower(day) {
		case "sunday":
			weekends = append(weekends, time.Sunday)
		case "monday":
			weekends = append(weekends, time.Monday)
		case "tuesday":
			weekends = append(weekends, time.Tuesday)
		case "wednesday":
			weekends = append(weekends, time.Wednesday)
		case "thursday":
			weekends = append(weekends, time.Thursday)
		case "friday":
			weekends = append(weekends, time.Friday)
		case "saturday":
			weekends = append(weekends, time.Saturday)
		}
	}

	return OutputSlackConfig{
		OfficeTimeStart:   input.OfficeTimeStart,
		OfficeTimeFinish:  input.OfficeTimeFinish,
		OfficeTimeWeekend: weekends,
	}
}
