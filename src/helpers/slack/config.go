package slack

import (
	helpers "dotfiles/src/helpers"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
)

type inputOfficeTimeOffDays struct {
	Jan []int
	Feb []int
	Mar []int
	Apr []int
	May []int
	Jun []int
	Jul []int
	Aug []int
	Sep []int
	Oct []int
	Nov []int
	Dec []int
}

type inputSlackConfig struct {
	OfficeTimeStart    int
	OfficeTimeFinish   int
	OfficeTimeWeekends []string
	OfficeTimeOffDays  inputOfficeTimeOffDays
}

type OutputSlackConfig struct {
	OfficeTimeStart    int
	OfficeTimeFinish   int
	OfficeTimeOffDays  []string // format: "month:day"
	OfficeTimeWeekends []time.Weekday
}

func ReadSlackConfig() OutputSlackConfig {
	configInput := helpers.ReadConfig[inputSlackConfig]("@/config/slack-status.jsonc")
	weekends := generateWeekends(configInput.OfficeTimeWeekends)
	offDays := generateOffDays(configInput.OfficeTimeOffDays)

	return OutputSlackConfig{
		OfficeTimeStart:    configInput.OfficeTimeStart,
		OfficeTimeFinish:   configInput.OfficeTimeFinish,
		OfficeTimeWeekends: weekends,
		OfficeTimeOffDays:  offDays,
	}
}

func GenerateOffDaysHash(month time.Month, day int) string {
	return month.String() + " " + strconv.Itoa(day)
}

func generateWeekends(input []string) []time.Weekday {
	weekends := make([]time.Weekday, 0, len(input))

	for _, day := range input {
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

	return weekends
}

func generateOffDays(input inputOfficeTimeOffDays) []string {
	offDays := make([]string, 0)

	for m := time.January; m <= time.December; m++ {
		shortName := m.String()[:3]

		field := reflect.ValueOf(input).FieldByName(shortName)
		if !field.IsValid() {
			fmt.Println(aurora.Red("Error: Invalid field name: " + shortName))
			continue
		}

		days, ok := field.Interface().([]int)
		if !ok {
			fmt.Println(aurora.Red("Error: Invalid field type: " + shortName))
			continue
		}

		for _, day := range days {
			offDays = append(offDays, GenerateOffDaysHash(m, day))
		}
	}

	return offDays
}
