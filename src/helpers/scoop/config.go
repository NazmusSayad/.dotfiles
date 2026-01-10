package scoop

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type ScoopAppInputConfig struct {
	ID            string
	Label         string
	Version       string
	SkipHashCheck bool
}

type ScoopAppConfig struct {
	Name          string
	Source        string
	Version       string
	SkipHashCheck bool

	Label string
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
			Name:   appId,
			Source: bucketId,

			Label:         app.Label,
			Version:       app.Version,
			SkipHashCheck: app.SkipHashCheck,
		})
	}

	return outputConfig
}

func GetScoopConfigAppMap(configs []ScoopAppConfig) map[string]ScoopAppConfig {
	appMap := make(map[string]ScoopAppConfig)

	for _, app := range configs {
		appMap[app.Source+"/"+app.Name] = app
	}

	return appMap
}

func GetScoopConfigBucketsList(configs []ScoopAppConfig) []string {
	bucketList := []string{}

	for _, app := range configs {
		bucketList = append(bucketList, app.Source)
	}

	return utils.UniqueArray(bucketList)
}
