package scoop

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"slices"
	"strconv"

	"github.com/logrusorgru/aurora/v4"
)

func SyncApps(configs []ScoopAppConfig, exports Export) {
	configApps := []string{}
	availableApps := []string{}
	unavailableApps := []ScoopAppConfig{}
	unnecessaryApps := []string{}

	for _, app := range configs {
		configApps = append(configApps, app.ID)
	}

	for _, app := range exports.Apps {
		availableApps = append(availableApps, app.Name)
	}

	for _, app := range configs {
		if !slices.Contains(availableApps, app.ID) {
			unavailableApps = append(unavailableApps, app)
		}
	}

	for _, app := range availableApps {
		if !slices.Contains(configApps, app) {
			unnecessaryApps = append(unnecessaryApps, app)
		}
	}

	fmt.Println(
		"Total: "+aurora.Green(strconv.Itoa(len(configApps))).String(),
		aurora.Faint("|").String()+" Available: "+aurora.Green(strconv.Itoa(len(availableApps))).String(),
		aurora.Faint("|").String()+" Unavailable: "+aurora.Green(strconv.Itoa(len(unavailableApps))).String(),
		aurora.Faint("|").String()+" Unnecessary: "+aurora.Green(strconv.Itoa(len(unnecessaryApps))).String(),
	)

	for _, app := range unavailableApps {
		fmt.Println(aurora.Faint("- Installing app ").String() + aurora.Green(app.Name).String())

		app := utils.Ternary(app.Version != "", app.ID+"@"+app.Version, app.ID)
		helpers.ExecNativeCommand([]string{"scoop", "install", app})

		fmt.Println()
	}

	for _, app := range unnecessaryApps {
		fmt.Println(aurora.Faint("- Removing app ").String() + aurora.Green(app).String())
		helpers.ExecNativeCommand([]string{"scoop", "uninstall", app})
		fmt.Println()
	}
}
