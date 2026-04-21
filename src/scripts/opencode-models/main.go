package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

func main() {
	providerConfigPath := helpers.ResolvePath("@/config/ai/opencode-providers.yaml")
	providerConfigs := helpers.ReadConfig[map[string]opencode.OpencodeProviderConfig](providerConfigPath)

	authConfigPath := helpers.ResolvePath("~/.local/share/opencode/auth.json")
	authConfig := helpers.ReadConfig[opencode.AuthConfig](authConfigPath)

	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
	fmt.Println(aurora.Cyan("Reading the OpenCode configuration...").String())
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("failed to read opencode config:", err)
		os.Exit(1)
	}

	var config map[string]any
	if err := json.Unmarshal(jsonc.ToJSON(configBytes), &config); err != nil {
		fmt.Println("failed to decode opencode config:", err)
		os.Exit(1)
	}

	providers, ok := config["provider"].(map[string]any)
	if !ok || providers == nil {
		providers = map[string]any{}
	}

	for providerID, providerConfig := range providerConfigs {
		fmt.Printf("%s %s\n", aurora.Blue("Syncing models for").String(), aurora.Bold(providerID).String())

		if providerConfig.ModelsURL == "" {
			fmt.Fprintf(os.Stderr, "%s provider %q has no apiURL/modelsURL, skipping\n", aurora.Yellow("warn:").String(), providerID)
			fmt.Println()
			continue
		}

		var providerAuth *opencode.AuthProvider
		if auth, ok := authConfig[providerID]; ok {
			providerAuth = &auth
		}

		models, err := opencode.FetchModels(providerConfig, providerAuth)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		}

		existingProvider, exists := providers[providerID]
		if !exists {
			providers[providerID] = map[string]any{"models": models}
			fmt.Println()
			continue
		}

		providerObject, ok := existingProvider.(map[string]any)
		if !ok {
			fmt.Fprintf(os.Stderr, "%s provider %q is not an object, skipping\n", aurora.Yellow("warn:").String(), providerID)
			fmt.Println()
			continue
		}

		providerObject["models"] = models
		providers[providerID] = providerObject
		fmt.Println()
	}

	config["provider"] = providers

	newConfigBytes, err := json.Marshal(config)
	if err != nil {
		fmt.Println("failed to encode config:", err)
		os.Exit(1)
	}

	mergedConfigRaw, err := helpers.MergeJSONObject(string(configBytes), string(newConfigBytes))
	if err != nil {
		fmt.Println("failed to merge config:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Writing the updated OpenCode configuration...").String())
	if err := os.WriteFile(configPath, []byte(mergedConfigRaw), 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Updated OpenCode configuration successfully.").String())
}
