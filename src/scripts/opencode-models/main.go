package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

func main() {
	providerConfigs := opencode.ReadOpencodeProvidersConfig()
	authConfigPath := helpers.ResolvePath("~/.local/share/opencode/auth.json")
	authConfig := helpers.ReadConfig[opencode.AuthConfig](authConfigPath)

	fmt.Println(aurora.Cyan("Reading the OpenCode configuration...").String())
	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("failed to read opencode config:", err)
		os.Exit(1)
	}

	var fullConfig map[string]any
	if err := json.Unmarshal(jsonc.ToJSON(configBytes), &fullConfig); err != nil {
		fmt.Println("failed to decode opencode config:", err)
		os.Exit(1)
	}

	modelsDotDevResponse, modelsDotDevError := opencode.FetchModelsDotDev()
	if modelsDotDevError != nil {
		fmt.Println("failed to fetch models.dev models:", modelsDotDevError)
		return
	}

	outputProviderConfig := make(map[string]opencode.OpencodeOutputProviderConfig)
	fmt.Println()

	for providerID, providerConfig := range providerConfigs {
		fmt.Printf("%s %s\n", aurora.Blue("Syncing models for").String(), aurora.Bold(providerID).String())

		devModels := modelsDotDevResponse[providerID]
		result, err := opencode.ResolveOpencodeProvider(providerID, providerConfig, devModels, authConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed to resolve provider %q: %v\n", aurora.Yellow("warn:").String(), providerID, err)
			fmt.Println()
			continue
		}

		outputProviderConfig[providerID] = result
		fmt.Println()
	}

	enabledProviders := make([]string, 0)
	for providerID := range outputProviderConfig {
		enabledProviders = append(enabledProviders, providerID)
	}

	fullConfig["provider"] = outputProviderConfig
	fullConfig["enabled_providers"] = utils.SortArray(enabledProviders)

	newConfigBytes, err := json.Marshal(fullConfig)
	if err != nil {
		fmt.Println("failed to encode config:", err)
		os.Exit(1)
	}

	mergedConfigRaw, err := helpers.MergeJSONObject(string(configBytes), string(newConfigBytes))
	if err != nil {
		fmt.Println("failed to merge config:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, []byte(mergedConfigRaw), 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	refreshErr := helpers.ExecNativeCommand([]string{"opencode", "models", "--refresh"})
	if refreshErr != nil {
		fmt.Println("failed to refresh opencode models")
	}
}
