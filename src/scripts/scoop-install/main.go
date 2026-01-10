package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/helpers/scoop"
	"fmt"
	"slices"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	exports := scoop.GetScoopExports()
	configs := scoop.ReadScoopAppConfig()

	configBucketsList := scoop.GetScoopConfigBucketsList(configs)
	exportBucketList := scoop.GetScoopExportBucketsList(exports)
	configAppMap := scoop.GetScoopConfigAppMap(configs)
	exportAppMap := scoop.GetScoopExportAppMap(exports)

	unavailableBuckets := []string{}
	for _, bucket := range configBucketsList {
		if !slices.Contains(exportBucketList, bucket) {
			unavailableBuckets = append(unavailableBuckets, bucket)
		}
	}

	unavailableApps := []scoop.ScoopAppConfig{}
	for appId, configApp := range configAppMap {
		_, isExists := exportAppMap[appId]

		if !isExists {
			unavailableApps = append(unavailableApps, configApp)
		}
	}

	unavailableBucketsCount := len(unavailableBuckets)
	if unavailableBucketsCount == 0 {
		fmt.Println(aurora.Green("No unavailable buckets found"))
	} else {
		fmt.Println("> Unavailable buckets:", aurora.Red(strings.Join(unavailableBuckets, ", ")))
	}

	unavailableAppsCount := len(unavailableApps)
	if unavailableAppsCount == 0 {
		fmt.Println(aurora.Green("No unavailable apps found"))
	} else {
		appNames := []string{}
		for _, app := range unavailableApps {
			appNames = append(appNames, app.Source+"/"+app.Name)
		}

		fmt.Println("> Unavailable apps:", aurora.Red(strings.Join(appNames, ", ")))
	}

	for _, bucket := range unavailableBuckets {
		fmt.Println()
		fmt.Println(aurora.Faint("- Installing bucket ").String() + aurora.Green(bucket).String())
		helpers.ExecNativeCommand([]string{"scoop", "bucket", "add", bucket})
	}

	for _, app := range unavailableApps {
		fmt.Println()
		fmt.Println(aurora.Faint("- Installing app ").String() + aurora.Green(app.Source+"/"+app.Name).String())
		helpers.ExecNativeCommand([]string{"scoop", "install", app.Source + "/" + app.Name})
	}
}
