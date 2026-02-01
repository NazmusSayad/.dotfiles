package scoop

import (
	"dotfiles/src/helpers"
	"fmt"
	"slices"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

func InstallScoopApps() {
	exports := GetScoopExports()
	configs := ReadScoopAppConfig()

	configBucketsList := GetScoopConfigBucketsList(configs)
	configAppMap := GetScoopConfigAppMap(configs)

	exportBucketList := GetScoopExportBucketsList(exports)
	exportAppMap := GetScoopExportAppMap(exports)

	missingBuckets := []string{}
	for _, bucket := range configBucketsList {
		if !slices.Contains(exportBucketList, bucket) {
			missingBuckets = append(missingBuckets, bucket)
		}
	}

	missingApps := []ScoopAppConfig{}
	for appId, configApp := range configAppMap {
		_, isExists := exportAppMap[appId]

		if !isExists {
			missingApps = append(missingApps, configApp)
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
		fmt.Println(aurora.Faint("- Installing bucket ").String() + aurora.Green(bucket).String())
		helpers.ExecNativeCommand([]string{"scoop", "bucket", "add", bucket})
	}

	for _, app := range missingApps {
		fmt.Println()
		fmt.Println(aurora.Faint("- Installing app ").String() + aurora.Green(app.ID).String())
		helpers.ExecNativeCommand([]string{"scoop", "install", app.ID})
	}
}
