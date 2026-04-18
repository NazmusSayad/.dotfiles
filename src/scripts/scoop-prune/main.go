package main

import (
	"fmt"
	"slices"
	"strings"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/scoop"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	configs := helpers.GetScoopApps()
	exports := scoop.GetScoopExports()

	configBucketsList := scoop.GetScoopConfigBucketsList(configs)
	configAppMap := scoop.GetScoopConfigAppMap(configs)
	configSrcMap := scoop.GetScoopConfigSrcMap(configs)

	exportBucketList := scoop.GetScoopExportBucketsList(exports)
	exportAppMap := scoop.GetScoopExportAppMap(exports)

	unnecessaryBuckets := []string{}
	for _, bucket := range exportBucketList {
		if !slices.Contains(configBucketsList, bucket) {
			unnecessaryBuckets = append(unnecessaryBuckets, bucket)
		}
	}

	unnecessaryApps := []scoop.ScoopApp{}
	for appId, exportApp := range exportAppMap {
		if slices.Contains(scoop.SCOOP_SYSTEM_APPS, appId) {
			continue
		}

		_, isExists := configAppMap[appId]
		if isExists {
			continue
		}

		if strings.HasSuffix(exportApp.Source, ".json") {
			_, isSourceExists := configSrcMap[exportApp.Source]
			if isSourceExists {
				continue
			}
		}

		unnecessaryApps = append(unnecessaryApps, exportApp)
	}

	unnecessaryBucketsCount := len(unnecessaryBuckets)
	if unnecessaryBucketsCount == 0 {
		fmt.Println(aurora.Green("No unnecessary buckets found"))
	} else {
		fmt.Println("> Unnecessary buckets:", aurora.Red(strings.Join(unnecessaryBuckets, ", ")))
	}

	unnecessaryAppsCount := len(unnecessaryApps)
	if unnecessaryAppsCount == 0 {
		fmt.Println(aurora.Green("No unnecessary apps found"))
	} else {
		appNames := []string{}
		for _, app := range unnecessaryApps {
			appNames = append(appNames, app.Source+"/"+app.Name)
		}

		fmt.Println("> Unnecessary apps:", aurora.Red(strings.Join(appNames, ", ")))
	}

	for _, bucket := range unnecessaryBuckets {
		fmt.Println()
		fmt.Println(aurora.Faint("- Removing bucket"), aurora.Green(bucket))
		helpers.ExecNativeCommand(
			[]string{"scoop", "bucket", "rm", bucket},
			helpers.ExecCommandOptions{Simulate: true},
		)
	}

	for _, app := range unnecessaryApps {
		fmt.Println()
		fmt.Println(aurora.Faint("- Removing app"), aurora.Green(app.Source+"/"+app.Name))

		if strings.HasSuffix(app.Source, ".json") {
			helpers.ExecNativeCommand(
				[]string{"scoop", "uninstall", app.Name},
				helpers.ExecCommandOptions{Simulate: true},
			)
		} else {
			helpers.ExecNativeCommand(
				[]string{"scoop", "uninstall", app.Source + "/" + app.Name},
				helpers.ExecCommandOptions{Simulate: true},
			)
		}
	}
}
