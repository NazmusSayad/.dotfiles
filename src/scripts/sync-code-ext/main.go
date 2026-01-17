package main

import (
	"dotfiles/src/helpers"
	"dotfiles/src/utils"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type ExtensionConfig struct {
	Id      string
	Version string
	Include []string
	Exclude []string
}

type InstalledExtension struct {
	Id      string
	Version string
}

var executables = []string{
	"code",
	"code-insiders",
	"antigravity",
	"cursor",
	"windsurf",
	"trae",
	"kiro",
}

func main() {
	configExts := helpers.ReadConfig[[]ExtensionConfig]("@/config/vscode/extensions.jsonc")
	configExtsMap := make(map[string]ExtensionConfig)
	for _, ext := range configExts {
		if len(ext.Include) > 0 && len(ext.Exclude) > 0 {
			fmt.Println(aurora.Faint(ext.Id), aurora.Red("Include and Exclude cannot be used together"))
			os.Exit(1)
		}

		configExtsMap[ext.Id] = ext
	}

	for _, executable := range executables {
		fmt.Println()
		fmt.Println(aurora.Bold(executable).String())

		err, installedExtIds := getCodeExtensions(executable)
		if err != nil {
			fmt.Println(aurora.Red("Failed to get extensions"))
			continue
		}

		installedExtMap := make(map[string]InstalledExtension)
		mismatchedExts := []ExtensionConfig{}
		missingExts := []ExtensionConfig{}
		extraExtIds := []string{}

		for _, ext := range installedExtIds {
			installedExtMap[ext.Id] = ext
		}

		for _, extension := range installedExtIds {
			configExt, isInConfig := configExtsMap[extension.Id]
			if !isInConfig || !isIncluded(configExt, executable) || isExcluded(configExt, executable) {
				extraExtIds = append(extraExtIds, extension.Id)
			}
		}

		for _, configExt := range configExts {
			installedExt, isInstalled := installedExtMap[configExt.Id]
			if isInstalled && !isVersionMatched(installedExt.Version, configExt.Version) {
				mismatchedExts = append(mismatchedExts, configExt)
			}
			if !isInstalled && !isExcluded(configExt, executable) && isIncluded(configExt, executable) {
				missingExts = append(missingExts, configExt)
			}
		}

		if len(extraExtIds) == 0 && len(missingExts) == 0 && len(mismatchedExts) == 0 {
			fmt.Println(aurora.Green("Already up to date"))
			continue
		}

		for _, ext := range mismatchedExts {
			fmt.Println(aurora.Faint("- Installing extension with version"), aurora.Green(ext.Id))
			helpers.ExecNativeCommand([]string{executable, "--install-extension", ext.Id + "@" + ext.Version})
		}

		for _, ext := range missingExts {
			fmt.Println(aurora.Faint("- Installing extension "), aurora.Green(ext.Id))
			helpers.ExecNativeCommand([]string{executable, "--install-extension", ext.Id})
		}

		for _, id := range extraExtIds {
			fmt.Println(aurora.Faint("- Uninstalling extension "), aurora.Red(id))
			helpers.ExecNativeCommand([]string{executable, "--uninstall-extension", id})
		}
	}
}

func getCodeExtensions(executable string) (error, []InstalledExtension) {
	cmd := exec.Command(executable, "--list-extensions", "--show-versions")
	output, err := cmd.Output()
	if err != nil {
		return err, nil
	}

	extensions := []InstalledExtension{}
	lines := strings.SplitSeq(string(output), "\n")

	for line := range lines {
		ext := strings.TrimSpace(line)
		if ext == "" {
			continue
		}

		parts := strings.SplitN(ext, "@", 2)
		extensions = append(extensions, InstalledExtension{
			Id:      parts[0],
			Version: utils.Ternary(len(parts) > 1, parts[1], ""),
		})
	}

	return nil, extensions

}

func isIncluded(configExt ExtensionConfig, executable string) bool {
	return utils.Ternary(len(configExt.Include) > 0, slices.Contains(configExt.Include, executable), true)
}

func isExcluded(configExt ExtensionConfig, executable string) bool {
	return len(configExt.Exclude) > 0 && slices.Contains(configExt.Exclude, executable)
}

func isVersionMatched(installedVersion, requiredVersion string) bool {
	if requiredVersion == "" {
		return true
	}

	return installedVersion == requiredVersion
}
