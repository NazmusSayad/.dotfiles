package scoop

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
)

var SCOOP_SYSTEM_APPS = []string{"main/7zip", "main/innounp"}

func GetScoopConfigAppMap(configs []helpers.ScoopAppConfig) map[string]helpers.ScoopAppConfig {
	appMap := make(map[string]helpers.ScoopAppConfig)

	for _, app := range configs {
		appMap[app.ID] = app
	}

	return appMap
}

func GetScoopConfigSrcMap(configs []helpers.ScoopAppConfig) map[string]helpers.ScoopAppConfig {
	appMap := make(map[string]helpers.ScoopAppConfig)

	for _, app := range configs {
		appMap[app.Source] = app
	}

	return appMap
}

func GetScoopConfigBucketsList(configs []helpers.ScoopAppConfig) []string {
	bucketList := []string{}

	for _, app := range configs {
		bucketList = append(bucketList, app.Bucket)
	}

	return utils.UniqueArray(bucketList)
}
