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

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

type authProvider struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type authConfig map[string]authProvider

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
	ID            string                       `json:"id"`
	Name          string                       `json:"name"`
	ContextLength int                          `json:"context_length"`
	Architecture  openAiCompatibleArchitecture `json:"architecture"`
	Opencode      kiloOptionalOpencode         `json:"opencode"`
}

type openAiCompatibleArchitecture struct {
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
}

type kiloOptionalOpencode struct {
	Family   string                     `json:"family"`
	Variants map[string]json.RawMessage `json:"variants"`
}

type opencodeOutputModel struct {
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	Limit      *opencodeOutputLimit       `json:"limit,omitempty"`
	Modalities *opencodeOutputModalities  `json:"modalities,omitempty"`
	Family     string                     `json:"family,omitempty"`
	Variants   map[string]json.RawMessage `json:"variants,omitempty"`
}

type opencodeOutputModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
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

var (
	managedProviderSuffix = "+"
	allowedModalities     = []string{"text", "audio", "image", "video", "pdf"}
)

func main() {
	providerConfigPath := helpers.ResolvePath("@/config/ai/opencode-providers.jsonc")
	providerConfigs := helpers.ReadConfig[map[string]opencodeProviderConfig](providerConfigPath)

	authConfigPath := helpers.ResolvePath("~/.local/share/opencode/auth.json")
	authConfig := helpers.ReadConfig[authConfig](authConfigPath)

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

	fmt.Println()

	desiredManagedProviders := map[string]opencodeOutputProvider{}
	for providerID, providerConfig := range providerConfigs {
		if !strings.HasSuffix(providerID, managedProviderSuffix) {
			providerID += managedProviderSuffix
		}

		fmt.Printf("%s %s\n", aurora.Blue("Syncing models for").String(), aurora.Bold(providerConfig.Name).String())

		var providerAuth *authProvider
		if auth, ok := authConfig[providerID]; ok {
			providerAuth = &auth
		}

		models, err := fetchModels(providerConfig, providerAuth)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
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

		fmt.Println()
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

	fmt.Println(aurora.Green("Writing the updated OpenCode configuration...").String())
	if err := os.WriteFile(configPath, updatedConfig, 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	fmt.Println(aurora.Green("Updated OpenCode configuration successfully.").String())
}

func fetchModels(providerConfig opencodeProviderConfig, auth *authProvider) (map[string]opencodeOutputModel, error) {
	modelsURL := providerConfig.ModelsURL
	if modelsURL == "" {
		modelsURL = strings.TrimRight(providerConfig.BaseURL, "/") + "/models"
	}

	fmt.Printf("%s %s\n", aurora.Yellow("Fetching models from").String(), aurora.Faint(modelsURL).String())

	req, err := http.NewRequest("GET", modelsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s models: %w", providerConfig.Name, err)
	}
	if auth != nil && auth.Type == "api" && auth.Key != "" {
		req.Header.Set("Authorization", "Bearer "+auth.Key)
	}

	resp, err := http.DefaultClient.Do(req)
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
	matchedModelIDs := map[string]bool{}
	for _, model := range payload.Data {
		if !slices.Contains(providerConfig.Models, model.ID) {
			continue
		}

		matchedModelIDs[model.ID] = true

		modelName := model.Name
		if modelName == "" {
			modelName = model.ID
		}

		entry := opencodeOutputModel{ID: model.ID, Name: modelName}
		if providerConfig.HasTurboMode && strings.LastIndex(model.ID, ":") <= strings.LastIndex(model.ID, "/") {
			entry.ID += ":nitro"
			entry.Name += " ⚡"
		}

		if model.ContextLength > 0 {
			entry.Limit = &opencodeOutputLimit{Context: model.ContextLength, Output: model.ContextLength}
		}

		supportedInputModalities := filterLLMModalities(model.Architecture.InputModalities)
		supportedOutputModalities := filterLLMModalities(model.Architecture.OutputModalities)

		if len(supportedInputModalities) > 0 && len(supportedOutputModalities) > 0 {
			entry.Modalities = &opencodeOutputModalities{
				Input:  supportedInputModalities,
				Output: supportedOutputModalities,
			}
		}

		if model.Opencode.Family != "" {
			entry.Family = model.Opencode.Family
		}

		if len(model.Opencode.Variants) > 0 {
			entry.Variants = model.Opencode.Variants
		}

		models[modelName] = entry
	}

	for _, modelID := range providerConfig.Models {
		if matchedModelIDs[modelID] {
			continue
		}

		fmt.Fprintf(os.Stderr, "%s model %q was not found for provider %q, using ID as name\n", aurora.Yellow("warn:").String(), modelID, providerConfig.Name)

		entry := opencodeOutputModel{ID: modelID, Name: modelID}
		models[modelID] = entry
	}

	return models, nil
}

func filterLLMModalities(modalities []string) []string {
	var filtered []string
	for _, m := range modalities {
		if slices.Contains(allowedModalities, m) {
			filtered = append(filtered, m)
		}
	}

	return filtered
}
