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
	configAppMap := scoop.GetScoopConfigAppMap(configs)

	exportBucketList := scoop.GetScoopExportBucketsList(exports)
	exportAppMap := scoop.GetScoopExportAppMap(exports)

	missingBuckets := []string{}
	for _, bucket := range configBucketsList {
		if !slices.Contains(exportBucketList, bucket) {
			missingBuckets = append(missingBuckets, bucket)
		}
	}

	missingApps := []scoop.ScoopAppConfig{}
	for appId, configApp := range configAppMap {
		_, isExists := exportAppMap[appId]

		if !isExists {
			missingApps = append(missingApps, configApp)
		}
	}

	for _, app := range scoop.SCOOP_SYSTEM_APPS {
		if _, isExists := exportAppMap[app]; !isExists {
			fmt.Println()
			fmt.Println(aurora.Red(app + " is required to install other apps"))
			installScoopApp(app)
		}
	}

	missingBucketsCount := len(missingBuckets)
	if missingBucketsCount == 0 {
		fmt.Println(aurora.Green("No missing buckets found"))
	} else {
		fmt.Println("> Missing buckets:", aurora.Red(strings.Join(missingBuckets, ", ")))
	}

	missingAppsCount := len(missingApps)
	if missingAppsCount == 0 {
		fmt.Println(aurora.Green("No missing apps found"))
	} else {
		appNames := []string{}
		for _, app := range missingApps {
			appNames = append(appNames, app.ID)
		}

		fmt.Println("> Missing apps:", aurora.Red(strings.Join(appNames, ", ")))
	}

	for _, bucket := range missingBuckets {
		fmt.Println()
		installScoopBucket(bucket)
	}

	for _, app := range missingApps {
		fmt.Println()
		installScoopApp(app.ID)
	}
}

func installScoopApp(appId string) {
	fmt.Println(aurora.Faint("- Installing"), aurora.Green(appId))
	helpers.ExecNativeCommand([]string{"scoop", "install", appId})
}

func installScoopBucket(bucket string) {
	fmt.Println(aurora.Faint("- Installing"), aurora.Green(bucket))
	helpers.ExecNativeCommand([]string{"scoop", "bucket", "add", bucket})
}
