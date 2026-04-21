package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

var managedProviderSuffix = "+"

func main() {
	providerConfigPath := helpers.ResolvePath("@/config/ai/opencode-providers.yaml")
	providerConfigs := helpers.ReadConfig[map[string]opencode.OpencodeProviderConfig](providerConfigPath)

	authConfigPath := helpers.ResolvePath("~/.local/share/opencode/auth.json")
	authConfig := helpers.ReadConfig[opencode.AuthConfig](authConfigPath)

	managedProviders := map[string]opencode.OpencodeOutputProvider{}

	for providerID, providerConfig := range providerConfigs {
		if !strings.HasSuffix(providerID, managedProviderSuffix) {
			providerID += managedProviderSuffix
		}

		fmt.Printf("%s %s\n", aurora.Blue("Syncing models for").String(), aurora.Bold(providerConfig.Name).String())

		var providerAuth *opencode.AuthProvider
		if auth, ok := authConfig[providerID]; ok {
			providerAuth = &auth
		}

		models := map[string]opencode.OpencodeOutputModel{}
		var err error

		if providerConfig.ModelBaseURL != "" {
			models, err = opencode.FetchUnknownModels(providerConfig, providerAuth)
		} else {
			models, err = opencode.FetchModels(providerConfig, providerAuth)
		}

		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		}

		modelsJSON, err := json.Marshal(models)
		if err != nil {
			fmt.Println("failed to encode models:", err)
			os.Exit(1)
		}

		managedProviders[providerID] = opencode.OpencodeOutputProvider{
			API:    providerConfig.BaseURL,
			Name:   providerConfig.Name,
			Models: json.RawMessage(modelsJSON),
		}

		fmt.Println()
	}

	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
	fmt.Println(aurora.Cyan("Reading the OpenCode configuration...").String())
	config, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("failed to read opencode config:", err)
		os.Exit(1)
	}

	configJSON := jsonc.ToJSON(config)

	var root map[string]json.RawMessage
	if err := json.Unmarshal(configJSON, &root); err != nil {
		fmt.Println("failed to decode opencode config:", err)
		os.Exit(1)
	}

	providerRaw, ok := root["provider"]
	if !ok {
		fmt.Println("failed to find provider block in opencode config")
		os.Exit(1)
	}

	var providers map[string]json.RawMessage
	if err := json.Unmarshal(providerRaw, &providers); err != nil {
		fmt.Println("failed to decode provider block:", err)
		os.Exit(1)
	}

	newProviders := map[string]json.RawMessage{}
	for providerID, providerRaw := range providers {
		if strings.HasSuffix(providerID, managedProviderSuffix) {
			continue
		}
		newProviders[providerID] = providerRaw
	}

	for providerID, provider := range managedProviders {
		providerRaw, err := json.Marshal(provider)
		if err != nil {
			fmt.Println("failed to encode provider:", err)
			os.Exit(1)
		}
		newProviders[providerID] = providerRaw
	}

	newProviderRaw, err := json.Marshal(newProviders)
	if err != nil {
		fmt.Println("failed to encode provider block:", err)
		os.Exit(1)
	}

	mergedProviderRaw, err := helpers.MergeJSONObject(string(providerRaw), string(newProviderRaw))
	if err != nil {
		fmt.Println("failed to merge provider block:", err)
		os.Exit(1)
	}

	providerIndex := strings.Index(string(config), string(providerRaw))
	if providerIndex == -1 {
		fmt.Println("failed to locate provider block in opencode config")
		os.Exit(1)
	}

	updatedConfig := append([]byte{}, config[:providerIndex]...)
	updatedConfig = append(updatedConfig, mergedProviderRaw...)
	updatedConfig = append(updatedConfig, config[providerIndex+len(providerRaw):]...)

	fmt.Println(aurora.Green("Writing the updated OpenCode configuration...").String())
	if err := os.WriteFile(configPath, updatedConfig, 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Updated OpenCode configuration successfully.").String())
}
