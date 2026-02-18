package scoop

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type ScoopAppInputConfig struct {
	ID    string
	Label string
}

type ScoopAppConfig struct {
	ID     string
	Bucket string
}

var SCOOP_SYSTEM_APPS = []string{"main/7zip", "main/git", "main/innounp"}

func ReadScoopAppConfig() []ScoopAppConfig {
	inputConfig := helpers.ReadConfig[[]string]("@/config/scoop-apps.yaml")
	outputConfig := []ScoopAppConfig{}

	for _, app := range inputConfig {
		appName := ""
		bucketName := ""
		splitStr := strings.Split(app, "/")

		if len(splitStr) == 1 {
			appName = splitStr[0]
			bucketName = "main"
		} else if len(splitStr) == 2 {
			bucketName = splitStr[0]
			appName = splitStr[1]
		} else {
			fmt.Println(aurora.Red("Invalid app ID; expected: <bucket>/<app>"))
			continue
		}

		outputConfig = append(outputConfig, ScoopAppConfig{
			ID:     bucketName + "/" + appName,
			Bucket: bucketName,
		})
	}

	return outputConfig
}

func GetScoopConfigAppMap(configs []ScoopAppConfig) map[string]ScoopAppConfig {
	appMap := make(map[string]ScoopAppConfig)

	for _, app := range configs {
		appMap[app.ID] = app
	}

	return appMap
}

func GetScoopConfigBucketsList(configs []ScoopAppConfig) []string {
	bucketList := []string{}

	for _, app := range configs {
		bucketList = append(bucketList, app.Bucket)
	}

	return utils.UniqueArray(bucketList)
}
