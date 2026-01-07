package scoop

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"slices"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

var SYSTEM_BUCKETS = []string{"main"}

func SyncBuckets(configs []ScoopAppConfig, exports Export) {
	configBuckets := []string{}
	availableBuckets := []string{}
	unavailableBuckets := []string{}
	unnecessaryBuckets := []string{}

	for _, app := range configs {
		configBuckets = append(configBuckets, app.Bucket)
	}

	configBuckets = utils.UniqueArray(configBuckets)

	for _, bk := range exports.Buckets {
		availableBuckets = append(availableBuckets, bk.Name)
	}

	for _, bk := range configBuckets {
		if !slices.Contains(availableBuckets, bk) {
			unavailableBuckets = append(unavailableBuckets, bk)
		}
	}

	for _, bk := range availableBuckets {
		if !slices.Contains(configBuckets, bk) && !slices.Contains(SYSTEM_BUCKETS, bk) {
			unnecessaryBuckets = append(unnecessaryBuckets, bk)
		}
	}

	fmt.Println(
		"Total: "+aurora.Green(strconv.Itoa(len(configBuckets))).String(),
		aurora.Faint("|").String()+" Available: "+aurora.Green(strconv.Itoa(len(availableBuckets))).String(),
		aurora.Faint("|").String()+" Unavailable: "+aurora.Green(strconv.Itoa(len(unavailableBuckets))).String(),
		aurora.Faint("|").String()+" Unnecessary: "+aurora.Green(strconv.Itoa(len(unnecessaryBuckets))).String(),
	)

	for _, bk := range unavailableBuckets {
		fmt.Println(aurora.Faint("- Installing bucket ").String() + aurora.Green(bk).String())
		helpers.ExecNativeCommand([]string{
			"scoop", "bucket", "add", bk,
		}, helpers.ExecCommandOptions{
			Silent: true,
		})
	}

	for _, bk := range unnecessaryBuckets {
		fmt.Println(aurora.Faint("- Removing bucket ").String() + aurora.Green(bk).String())
		helpers.ExecNativeCommand([]string{
			"scoop", "bucket", "rm", bk,
		}, helpers.ExecCommandOptions{
			Silent: true,
		})
	}
}
