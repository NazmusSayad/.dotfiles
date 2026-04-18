package helpers

import (
	"fmt"
	"strings"

	"dotfiles/src/constants"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

type MsysAppConfig struct {
	ID   string `yaml:"id"`
	Repo string `yaml:"repo"`
}

type ScoopAppConfig struct {
	ID     string `yaml:"id"`
	Bucket string `yaml:"bucket"`
	Source string `yaml:"source"`
}

type WingetAppConfig struct {
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Scope         string `yaml:"scope"`
	Version       string `yaml:"version"`
	InstallerType string `yaml:"installerType"`

	ForceAdminInstall bool `yaml:"forceAdminInstall"`
	ForceAdminUpgrade bool `yaml:"forceAdminUpgrade"`

	InteractiveInstall bool `yaml:"interactiveInstall"`
	InteractiveUpgrade bool `yaml:"interactiveUpgrade"`

	SkipInstall      bool `yaml:"skipInstall"`
	SkipUpgrade      bool `yaml:"skipUpgrade"`
	SkipDependencies bool `yaml:"skipDependencies"`
}

func GetMsysApps() []MsysAppConfig {
	config := ReadConfig[map[string][]any]("@/config/apps.yaml")
	apps := []MsysAppConfig{}

	for _, item := range config["msys2"] {
		strApp, isString := item.(string)
		if isString {
			apps = append(apps, MsysAppConfig{ID: strApp})
			continue
		}

		objApp, isMap := item.(map[string]any)
		if isMap {
			id, _ := objApp["id"].(string)
			repo, _ := objApp["repo"].(string)
			apps = append(apps, MsysAppConfig{ID: id, Repo: repo})
		}
	}

	return apps
}

func GetScoopApps() []ScoopAppConfig {
	config := ReadConfig[map[string][]any]("@/config/apps.yaml")
	outputConfig := []ScoopAppConfig{}

	for _, item := range config["scoop"] {
		app := ""
		configuredBucket := ""

		s, isString := item.(string)
		if isString {
			app = s
		} else {
			m, isMap := item.(map[string]any)
			if !isMap {
				continue
			}

			app, _ = m["id"].(string)
			configuredBucket, _ = m["bucket"].(string)
		}

		appName := ""
		appSource := ""
		bucketName := ""
		splitStr := strings.Split(app, "/")

		if len(splitStr) == 1 {
			appName = splitStr[0]

			if strings.HasPrefix(appName, "$") {
				bucketName = ""
				appName = appName[1:]
				appSource = ResolvePath("@/" + constants.SCOOP_DIR + "/" + appName + ".json")
			} else {
				bucketName = utils.Ternary(configuredBucket != "", configuredBucket, "main")
			}
		} else if len(splitStr) == 2 {
			if configuredBucket != "" {
				fmt.Println(aurora.Red("Invalid app and bucket configuration; expected: <bucket>/<app>"))
				continue
			}

			bucketName = splitStr[0]
			appName = splitStr[1]
		} else {
			fmt.Println(aurora.Red("Invalid app ID; expected: <bucket>/<app>"))
			continue
		}

		outputConfig = append(outputConfig, ScoopAppConfig{
			ID:     utils.Ternary(bucketName == "", appName, bucketName+"/"+appName),
			Source: appSource,
			Bucket: bucketName,
		})
	}

	return outputConfig
}

func GetWingetApps() []WingetAppConfig {
	config := ReadConfig[map[string][]any]("@/config/apps.yaml")
	apps := []WingetAppConfig{}

	for _, item := range config["winget"] {
		s, isString := item.(string)
		if isString {
			apps = append(apps, WingetAppConfig{ID: s})
			continue
		}

		m, isMap := item.(map[string]any)
		if !isMap {
			continue
		}

		app := WingetAppConfig{}
		app.ID, _ = m["id"].(string)
		app.Name, _ = m["name"].(string)
		app.Scope, _ = m["scope"].(string)
		app.Version, _ = m["version"].(string)
		app.InstallerType, _ = m["installerType"].(string)
		app.ForceAdminInstall, _ = m["forceAdminInstall"].(bool)
		app.ForceAdminUpgrade, _ = m["forceAdminUpgrade"].(bool)
		app.InteractiveInstall, _ = m["interactiveInstall"].(bool)
		app.InteractiveUpgrade, _ = m["interactiveUpgrade"].(bool)
		app.SkipInstall, _ = m["skipInstall"].(bool)
		app.SkipUpgrade, _ = m["skipUpgrade"].(bool)
		app.SkipDependencies, _ = m["skipDependencies"].(bool)

		apps = append(apps, app)
	}

	return apps
}
