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

	modelsByProvider := map[string]json.RawMessage{}

	for providerID, providerConfig := range providerConfigs {
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

		modelsByProvider[providerID] = modelsJSON

		fmt.Println()
	}

	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
	fmt.Println(aurora.Cyan("Reading the OpenCode configuration...").String())
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("failed to read opencode config:", err)
		os.Exit(1)
	}

	oldConfigStr := string(configBytes)

	var root map[string]json.RawMessage
	if err := json.Unmarshal(jsonc.ToJSON(configBytes), &root); err != nil {
		fmt.Println("failed to decode opencode config:", err)
		os.Exit(1)
	}

	var providers map[string]json.RawMessage
	if err := json.Unmarshal(root["provider"], &providers); err != nil {
		fmt.Println("failed to decode provider block:", err)
		os.Exit(1)
	}

	for providerID, modelsJSON := range modelsByProvider {
		providerRaw, ok := providers[providerID]
		if !ok {
			fmt.Fprintf(os.Stderr, "%s provider %q not found in config, skipping\n", aurora.Yellow("warn:").String(), providerID)
			continue
		}

		var provider map[string]json.RawMessage
		if err := json.Unmarshal(providerRaw, &provider); err != nil {
			fmt.Fprintf(os.Stderr, "%s failed to decode provider %q: %v\n", aurora.Yellow("warn:").String(), providerID, err)
			continue
		}

		provider["models"] = modelsJSON

		updatedProviderRaw, err := json.Marshal(provider)
		if err != nil {
			fmt.Println("failed to encode provider:", err)
			os.Exit(1)
		}

		providers[providerID] = updatedProviderRaw
	}

	newProviderBlock, err := json.Marshal(providers)
	if err != nil {
		fmt.Println("failed to encode provider block:", err)
		os.Exit(1)
	}

	root["provider"] = newProviderBlock

	newConfigBytes, err := json.Marshal(root)
	if err != nil {
		fmt.Println("failed to encode config:", err)
		os.Exit(1)
	}

	mergedConfigStr, err := helpers.MergeJSONObject(oldConfigStr, string(newConfigBytes))
	if err != nil {
		fmt.Println("failed to merge config:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Writing the updated OpenCode configuration...").String())
	if err := os.WriteFile(configPath, []byte(mergedConfigStr), 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Updated OpenCode configuration successfully.").String())
}
