package scoop

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
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

	missingApps := []ScoopConfig{}
	for _, config := range configs {
		if !slices.Contains(installedApps, config.ID) {
			missingApps = append(missingApps, config)
		}
	}

	unnecessaryApps := []string{}
	for _, installedApp := range installedApps {
		if !slices.Contains(apps, installedApp) {
			unnecessaryApps = append(unnecessaryApps, installedApp)
		}
	}

	fmt.Println(
		"Total: "+aurora.Green(strconv.Itoa(len(apps))).String(),
		aurora.Faint("|").String()+" Installed: "+aurora.Green(strconv.Itoa(len(installedApps))).String(),
		aurora.Faint("|").String()+" Missing: "+aurora.Green(strconv.Itoa(len(missingApps))).String(),
		aurora.Faint("|").String()+" Extra: "+aurora.Green(strconv.Itoa(len(unnecessaryApps))).String(),
	)

	for _, missingApp := range missingApps {
		fmt.Println(aurora.Faint("- Installing app ").String() + aurora.Green(missingApp.Name).String())

		app := utils.Ternary(missingApp.Version != "", missingApp.ID+"@"+missingApp.Version, missingApp.ID)
		helpers.ExecNativeCommand([]string{"scoop", "install", app})

		fmt.Println()
	}

	for _, unnecessaryApp := range unnecessaryApps {
		fmt.Println(aurora.Faint("- Removing app ").String() + aurora.Green(unnecessaryApp).String())
		helpers.ExecNativeCommand([]string{"scoop", "uninstall", unnecessaryApp})
		fmt.Println()
	}
}
