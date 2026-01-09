package scoop

import (
	"dotfiles/src/helpers"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type ScoopAppInputConfig struct {
	ID            string
	Name          string
	Version       string
	SkipHashCheck bool
}

type ScoopAppConfig struct {
	ID            string
	Name          string
	Bucket        string
	Version       string
	SkipHashCheck bool
}

func ReadScoopAppConfig() []ScoopAppConfig {
	inputConfig := helpers.ReadConfig[[]ScoopAppInputConfig]("@/config/scoop-apps.jsonc")
	outputConfig := []ScoopAppConfig{}

	for _, app := range inputConfig {
		splitStr := strings.Split(app.ID, "/")
		bucketId := ""
		appId := ""

		if len(splitStr) == 0 {
			appId = app.ID
			bucketId = "main"
		} else if len(splitStr) == 1 {
			appId = app.ID
			bucketId = "main"
		} else if len(splitStr) == 2 {
			bucketId = splitStr[0]
			appId = splitStr[1]
		} else {
			fmt.Println(aurora.Red("Invalid app ID; expected: <bucket>/<app>"))
			continue
		}

		outputConfig = append(outputConfig, ScoopAppConfig{
			ID:     appId,
			Bucket: bucketId,

			Name:          app.Name,
			Version:       app.Version,
			SkipHashCheck: app.SkipHashCheck,
		})
	}

	return outputConfig
}
