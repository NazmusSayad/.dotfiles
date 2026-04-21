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
	providerConfigs := opencode.ReadOpencodeProvidersConfig()

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

	modelsDotDevResponse, modelsDotDevError := opencode.FetchModelsDotDev()
	if modelsDotDevError != nil {
		fmt.Println("failed to fetch models.dev models:", modelsDotDevError)
	}

	for providerID, providerConfig := range providerConfigs {
		fmt.Printf("%s %s\n", aurora.Blue("Syncing models for").String(), aurora.Bold(providerID).String())
		providerModelIDs := make([]string, 0, len(providerConfig.Models))
		for _, configuredModel := range providerConfig.Models {
			providerModelIDs = append(providerModelIDs, configuredModel.ID)
		}

		if providerConfig.ModelsURL == "" {
			devModels, ok := modelsDotDevResponse[providerID]
			if !ok {
				fmt.Fprintf(os.Stderr, "%s provider %q has no apiURL/modelsURL and no models.dev data, skipping\n", aurora.Yellow("warn:").String(), providerID)
				fmt.Println()
				continue
			}

			filteredModels := map[string]opencode.OpencodeOutputModel{}
			for _, configuredModel := range providerConfig.Models {
				modelID := configuredModel.ID
				model, ok := devModels[modelID]
				if !ok {
					fmt.Fprintf(os.Stderr, "%s model %q was not found for provider %q in models.dev, using ID as name\n", aurora.Yellow("warn:").String(), modelID, providerID)
					filteredModels[modelID] = opencode.ApplyModelContextCap(opencode.OpencodeOutputModel{ID: modelID, Name: modelID}, configuredModel.ContextCap)
					continue
				}
				model.ID = modelID
				filteredModels[modelID] = opencode.ApplyModelContextCap(model, configuredModel.ContextCap)
			}

			existingProvider, exists := providers[providerID]
			if !exists {
				providers[providerID] = map[string]any{
					"models":    filteredModels,
					"whitelist": providerModelIDs,
				}

				fmt.Println()
				continue
			}

			providerObject, ok := existingProvider.(map[string]any)
			if !ok {
				fmt.Fprintf(os.Stderr, "%s provider %q is not an object, skipping\n", aurora.Yellow("warn:").String(), providerID)
				fmt.Println()
				continue
			}

			providerObject["models"] = filteredModels
			providerObject["whitelist"] = providerModelIDs
			providers[providerID] = providerObject
			fmt.Println()
			continue
		}

		var providerAuth *opencode.AuthProvider
		if auth, ok := authConfig[providerID]; ok {
			providerAuth = &auth
		}

		models, err := opencode.FetchModels(providerID, providerConfig, providerAuth)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		}

		cappedModels := make(map[string]opencode.OpencodeOutputModel, len(providerConfig.Models))
		for _, configuredModel := range providerConfig.Models {
			modelID := configuredModel.ID
			model, ok := models[modelID]
			if !ok {
				model = opencode.OpencodeOutputModel{ID: modelID, Name: modelID}
			}

			cappedModels[modelID] = opencode.ApplyModelContextCap(model, configuredModel.ContextCap)
		}

		existingProvider, exists := providers[providerID]
		if !exists {
			providers[providerID] = map[string]any{"models": cappedModels, "whitelist": providerModelIDs}
			fmt.Println()
			continue
		}

		providerObject, ok := existingProvider.(map[string]any)
		if !ok {
			fmt.Fprintf(os.Stderr, "%s provider %q is not an object, skipping\n", aurora.Yellow("warn:").String(), providerID)
			fmt.Println()
			continue
		}

		providerObject["models"] = cappedModels
		providerObject["whitelist"] = providerModelIDs
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
