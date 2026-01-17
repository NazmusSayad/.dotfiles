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
	Include []string
	Exclude []string
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

		missingExts := []ExtensionConfig{}
		extraExtIds := []string{}

		for _, extension := range installedExtIds {
			configExt, isInConfig := configExtsMap[extension]
			if !isInConfig || !isIncluded(configExt, executable) || isExcluded(configExt, executable) {
				extraExtIds = append(extraExtIds, extension)
			}
		}

		for _, configExt := range configExts {
			isInstalled := slices.Contains(installedExtIds, configExt.Id)
			if !isInstalled && !isExcluded(configExt, executable) && isIncluded(configExt, executable) {
				missingExts = append(missingExts, configExt)
			}
		}

		if len(extraExtIds) == 0 && len(missingExts) == 0 {
			fmt.Println(aurora.Green("Already up to date"))
			continue
		}

		for _, ext := range missingExts {
			fmt.Println(aurora.Faint("- Installing extension ").String() + aurora.Green(ext.Id).String())
			helpers.ExecNativeCommand([]string{executable, "--install-extension", ext.Id})
		}

		for _, id := range extraExtIds {
			fmt.Println(aurora.Faint("- Uninstalling extension ").String() + aurora.Red(id).String())
			helpers.ExecNativeCommand([]string{executable, "--uninstall-extension", id})
		}
	}
}

func getCodeExtensions(executable string) (error, []string) {
	cmd := exec.Command(executable, "--list-extensions")
	output, err := cmd.Output()
	if err != nil {
		return err, nil
	}

	extensions := []string{}
	lines := strings.SplitSeq(string(output), "\n")

	for line := range lines {
		ext := strings.TrimSpace(line)
		if ext == "" {
			continue
		}

		extensions = append(extensions, ext)
	}

	return nil, extensions

}

func isIncluded(configExt ExtensionConfig, executable string) bool {
	return utils.Ternary(len(configExt.Include) > 0, slices.Contains(configExt.Include, executable), true)
}

func isExcluded(configExt ExtensionConfig, executable string) bool {
	return len(configExt.Exclude) > 0 && slices.Contains(configExt.Exclude, executable)
}
