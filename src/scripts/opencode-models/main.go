package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"dotfiles/src/helpers"

	"github.com/tidwall/jsonc"
)

type opencodeProviderConfig struct {
	Name         string   `json:"name"`
	BaseURL      string   `json:"apiURL"`
	ModelsURL    string   `json:"modelsURL"`
	HasTurboMode bool     `json:"hasTurboMode"`
	Models       []string `json:"models"`
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
	Name   string          `json:"name,omitempty"`
	Models json.RawMessage `json:"models"`
}

const managedProviderSuffix = "*"

func main() {
	providerConfigPath := helpers.ResolvePath("@/config/ai/opencode-providers.jsonc")
	providerConfigs := helpers.ReadConfig[map[string]opencodeProviderConfig](providerConfigPath)

	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
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
	for providerID, providerConfig := range providerConfigs {
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
			API:    providerConfig.BaseURL,
			Name:   providerConfig.Name,
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
		return nil, fmt.Errorf("failed to fetch %s models: %w", providerConfig.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s models: %s", providerConfig.Name, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s response body: %w", providerConfig.Name, err)
	}

	var payload openAiCompatibleModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", providerConfig.Name, err)
	}

	models := map[string]opencodeOutputModel{}
	for _, model := range payload.Data {
		if !slices.Contains(providerConfig.Models, model.ID) {
			continue
		}

		entry := opencodeOutputModel{ID: model.ID, Name: model.Name}
		if providerConfig.HasTurboMode && strings.LastIndex(model.ID, ":") <= strings.LastIndex(model.ID, "/") {
			entry.ID += ":turbo"
		}

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
