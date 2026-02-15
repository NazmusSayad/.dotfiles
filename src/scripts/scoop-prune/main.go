package scoop

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
	configAppMap := scoop.GetScoopConfigAppMap(configs)

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

		if _, isExists := configAppMap[appId]; !isExists {
			unnecessaryApps = append(unnecessaryApps, exportApp)
		}
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
		helpers.ExecNativeCommand(
			[]string{"scoop", "uninstall", app.Source + "/" + app.Name},
			helpers.ExecCommandOptions{Simulate: true},
		)
	}
}
