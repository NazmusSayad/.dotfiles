package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"dotfiles/src/helpers"

	"github.com/tidwall/jsonc"
)

type opencodeProviderConfig struct {
	ID        string   `json:"id"`
	APIURL    string   `json:"apiURL"`
	ModelsURL string   `json:"modelsURL"`
	Models    []string `json:"models"`
}

type openAiCompatibleModelsResponse struct {
	Data []openAiCompatibleModel `json:"data"`
}

type openAiCompatibleModel struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	ContextLength int                  `json:"context_length"`
	Opencode      kiloOptionalOpencode `json:"opencode"`
}

type kiloOptionalOpencode struct {
	Family   string                     `json:"family"`
	Variants map[string]json.RawMessage `json:"variants"`
}

type opencodeOutputModel struct {
	ID       string                     `json:"id"`
	Name     string                     `json:"name"`
	Limit    *opencodeOutputLimit       `json:"limit,omitempty"`
	Family   string                     `json:"family,omitempty"`
	Variants map[string]json.RawMessage `json:"variants,omitempty"`
}

type opencodeOutputLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

type opencodeOutputProvider struct {
	API    string          `json:"api,omitempty"`
	Models json.RawMessage `json:"models"`
}

const managedProviderSuffix = "+"

func main() {
	providerConfigPath := helpers.ResolvePath("@/config/ai/opencode-providers.jsonc")
	providerConfigContent, err := os.ReadFile(providerConfigPath)
	if err != nil {
		fmt.Println("failed to read provider config:", err)
		os.Exit(1)
	}

	var providerConfigs []opencodeProviderConfig
	if err := json.Unmarshal(jsonc.ToJSON(providerConfigContent), &providerConfigs); err != nil {
		fmt.Println("failed to decode provider config:", err)
		os.Exit(1)
	}

	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
	config, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("failed to read opencode config:", err)
		os.Exit(1)
	}

	var root map[string]json.RawMessage
	if err := json.Unmarshal(config, &root); err != nil {
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

	currentManagedProviders := map[string]opencodeOutputProvider{}
	for providerID, providerRaw := range providers {
		if strings.HasSuffix(providerID, managedProviderSuffix) {
			var provider opencodeOutputProvider
			if err := json.Unmarshal(providerRaw, &provider); err != nil {
				fmt.Println("failed to decode managed provider:", err)
				os.Exit(1)
			}

			currentManagedProviders[providerID] = provider
		}
	}

	desiredManagedProviders := map[string]opencodeOutputProvider{}
	for _, providerConfig := range providerConfigs {
		providerID := providerConfig.ID
		if !strings.HasSuffix(providerID, managedProviderSuffix) {
			providerID += managedProviderSuffix
		}

		models, err := fetchModels(providerConfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		desiredModels, err := json.Marshal(models)
		if err != nil {
			fmt.Println("failed to encode models:", err)
			os.Exit(1)
		}

		existingModels := "{}"
		if existing, ok := currentManagedProviders[providerID]; ok && len(existing.Models) > 0 {
			existingModels = string(existing.Models)
		}

		patchedModels, err := helpers.MergeJSONObject(existingModels, string(desiredModels))
		if err != nil {
			fmt.Println("failed to apply models patch:", err)
			os.Exit(1)
		}

		desiredManagedProviders[providerID] = opencodeOutputProvider{
			API:    providerConfig.APIURL,
			Models: json.RawMessage(patchedModels),
		}
	}

	finalProviders := map[string]json.RawMessage{}
	for providerID, providerRaw := range providers {
		if strings.HasSuffix(providerID, managedProviderSuffix) {
			continue
		}

		finalProviders[providerID] = providerRaw
	}

	for providerID, provider := range desiredManagedProviders {
		providerRaw, err := json.Marshal(provider)
		if err != nil {
			fmt.Println("failed to encode provider:", err)
			os.Exit(1)
		}

		finalProviders[providerID] = providerRaw
	}

	finalProviderRaw, err := json.Marshal(finalProviders)
	if err != nil {
		fmt.Println("failed to encode provider block:", err)
		os.Exit(1)
	}

	updatedProviderRaw, err := helpers.MergeJSONObject(string(providerRaw), string(finalProviderRaw))
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
	updatedConfig = append(updatedConfig, updatedProviderRaw...)
	updatedConfig = append(updatedConfig, config[providerIndex+len(providerRaw):]...)

	if err := os.WriteFile(configPath, updatedConfig, 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	fmt.Println("updated", configPath)
}

func fetchModels(providerConfig opencodeProviderConfig) (map[string]opencodeOutputModel, error) {
	resp, err := http.Get(providerConfig.ModelsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s models: %w", providerConfig.ID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s models: %s", providerConfig.ID, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s response body: %w", providerConfig.ID, err)
	}

	var payload openAiCompatibleModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", providerConfig.ID, err)
	}

	models := map[string]opencodeOutputModel{}
	for _, model := range payload.Data {
		if !contains(providerConfig.Models, model.ID) {
			continue
		}

		entry := opencodeOutputModel{ID: model.ID, Name: model.Name}
		if model.ContextLength > 0 {
			entry.Limit = &opencodeOutputLimit{Context: model.ContextLength, Output: model.ContextLength}
		}

		if model.Opencode.Family != "" {
			entry.Family = model.Opencode.Family
		}

		if len(model.Opencode.Variants) > 0 {
			entry.Variants = model.Opencode.Variants
		}

		models[model.Name] = entry
	}

	return models, nil
}

func contains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}

	return false
}
