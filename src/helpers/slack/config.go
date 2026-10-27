package slack

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	helpers "dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

type inputOfficeTimeOffDays struct {
	Jan []int `yaml:"Jan"`
	Feb []int `yaml:"Feb"`
	Mar []int `yaml:"Mar"`
	Apr []int `yaml:"Apr"`
	May []int `yaml:"May"`
	Jun []int `yaml:"Jun"`
	Jul []int `yaml:"Jul"`
	Aug []int `yaml:"Aug"`
	Sep []int `yaml:"Sep"`
	Oct []int `yaml:"Oct"`
	Nov []int `yaml:"Nov"`
	Dec []int `yaml:"Dec"`
}

type inputSlackConfig struct {
	OfficeTimeStart    int                    `yaml:"OfficeTimeStart"`
	OfficeTimeFinish   int                    `yaml:"OfficeTimeFinish"`
	OfficeTimeWeekends []string               `yaml:"OfficeTimeWeekends"`
	OfficeTimeOffDays  inputOfficeTimeOffDays `yaml:"OfficeTimeOffDays"`
}

type OutputSlackConfig struct {
	OfficeTimeStart    int
	OfficeTimeFinish   int
	OfficeTimeOffDays  []string // format: "month:day"
	OfficeTimeWeekends []time.Weekday
}

func ReadSlackConfig() map[string]OutputSlackConfig {
	configInput := helpers.ReadConfig[map[string]inputSlackConfig](
		"$SYSTEM_FILE_DIR/slack-manager.yml",
		helpers.ReadConfigOptions{Silent: true},
	)

	configs := make(map[string]OutputSlackConfig, len(configInput))
	for workspace, input := range configInput {
		configs[workspace] = OutputSlackConfig{
			OfficeTimeStart:    input.OfficeTimeStart,
			OfficeTimeFinish:   input.OfficeTimeFinish,
			OfficeTimeWeekends: generateWeekends(input.OfficeTimeWeekends),
			OfficeTimeOffDays:  generateOffDays(input.OfficeTimeOffDays),
		}
	}

	return configs
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
