package scoop

import (
	"dotfiles/src/helpers"
	"fmt"
	"slices"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

var SYSTEM_BUCKETS = []string{"main"}

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

	unnecessaryBuckets := []string{}
	for _, installedBucket := range installedBuckets {
		if !slices.Contains(configBuckets, installedBucket) && !slices.Contains(SYSTEM_BUCKETS, installedBucket) {
			unnecessaryBuckets = append(unnecessaryBuckets, installedBucket)
		}
	}

	fmt.Println(
		"Total: "+aurora.Green(strconv.Itoa(len(configBuckets))).String(),
		aurora.Faint("|").String()+" Installed: "+aurora.Green(strconv.Itoa(len(installedBuckets))).String(),
		aurora.Faint("|").String()+" Missing: "+aurora.Green(strconv.Itoa(len(missingBuckets))).String(),
		aurora.Faint("|").String()+" Extra: "+aurora.Green(strconv.Itoa(len(unnecessaryBuckets))).String(),
	)

	for _, missingBucket := range missingBuckets {
		fmt.Println(aurora.Faint("- Installing bucket ").String() + aurora.Green(missingBucket).String())
		helpers.ExecNativeCommand([]string{
			"scoop", "bucket", "add", missingBucket,
		}, helpers.ExecCommandOptions{
			Silent: true,
		})
	}

	for _, unnecessaryBucket := range unnecessaryBuckets {
		fmt.Println(aurora.Faint("- Removing bucket ").String() + aurora.Green(unnecessaryBucket).String())
		helpers.ExecNativeCommand([]string{
			"scoop", "bucket", "rm", unnecessaryBucket,
		}, helpers.ExecCommandOptions{
			Silent: true,
		})
	}
}
