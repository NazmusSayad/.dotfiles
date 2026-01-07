package scoop

import (
	"dotfiles/src/helpers"
	"fmt"
	"slices"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

func SyncApps(configs []ScoopConfig, exports Export) {
	apps := []string{}
	for _, config := range configs {
		apps = append(apps, config.ID)
	}

	installedApps := []string{}
	for _, app := range exports.Apps {
		installedApps = append(installedApps, app.Name)
	}

	missingApps := []string{}
	for _, config := range configs {
		if !slices.Contains(installedApps, config.ID) {
			missingApps = append(missingApps, config.ID)
		}
	}

	fmt.Println(aurora.Faint("Syncing apps, total: ").String() + aurora.Green(strconv.Itoa(len(missingApps))).String())

	for _, missingApp := range missingApps {
		fmt.Println(aurora.Faint("- Installing app ").String() + aurora.Green(missingApp).String())
		err := helpers.ExecNativeCommand([]string{"scoop", "install", missingApp}, helpers.ExecCommandOptions{
			Silent: true,
		})

		if err != nil {
			fmt.Println(aurora.Red("Error installing app: " + err.Error()))
		}
	}
}
