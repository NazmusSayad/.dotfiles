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

	if _, isExists := exportAppMap[GIT_APP_ID]; !isExists {
		fmt.Println()
		fmt.Println(aurora.Red("Git is required to install other apps"))
		installScoopApp(GIT_APP_ID)
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
	helpers.ExecNativeCommand(
		[]string{"scoop", "install", appId},
		helpers.ExecCommandOptions{Simulate: true},
	)
}

func installScoopBucket(bucket string) {
	fmt.Println(aurora.Faint("- Installing"), aurora.Green(bucket))
	helpers.ExecNativeCommand(
		[]string{"scoop", "bucket", "add", bucket},
		helpers.ExecCommandOptions{Simulate: true},
	)
}
