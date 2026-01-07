package scoop

import (
	"dotfiles/src/helpers"
	"fmt"
	"slices"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

func SyncBuckets(configs []ScoopConfig, exports Export) {
	configBuckets := []string{}
	for _, config := range configs {
		configBuckets = append(configBuckets, config.Bucket)
	}

	installedBuckets := []string{}
	for _, bucket := range exports.Buckets {
		installedBuckets = append(installedBuckets, bucket.Name)
	}

	missingBuckets := []string{}
	for _, configBucket := range configBuckets {
		if !slices.Contains(installedBuckets, configBucket) {
			missingBuckets = append(missingBuckets, configBucket)
		}
	}

	fmt.Println(aurora.Faint("Syncing buckets, total: ").String() + aurora.Green(strconv.Itoa(len(missingBuckets))).String())

	for _, missingBucket := range missingBuckets {
		fmt.Println(aurora.Faint("- Installing bucket ").String() + aurora.Green(missingBucket).String())
		err := helpers.ExecNativeCommand([]string{
			"scoop", "bucket", "add", missingBucket,
		})

		if err != nil {
			fmt.Println(aurora.Red("Error installing bucket: " + err.Error()))
		} else {
			fmt.Println(aurora.Green("Bucket installed successfully"))
		}
	}
}
