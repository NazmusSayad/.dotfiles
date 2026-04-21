package opencode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

var allowedModalities = []string{"text", "audio", "image", "video", "pdf"}

func FetchModels(providerConfig OpencodeProviderConfig, auth *AuthProvider) (map[string]OpencodeOutputModel, error) {
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

	var payload OpenAiCompatibleModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", providerConfig.Name, err)
	}

	models := map[string]OpencodeOutputModel{}
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

		entry := OpencodeOutputModel{ID: model.ID, Name: modelName}
		if providerConfig.HasTurboMode && strings.LastIndex(model.ID, ":") <= strings.LastIndex(model.ID, "/") {
			entry.ID += ":nitro"
			entry.Name += " ⚡"
		}

		if providerConfig.ModelPrefix != "" {
			entry.Name = providerConfig.ModelPrefix + entry.Name
		}

		if model.ContextLength > 0 {
			entry.Limit = &OpencodeOutputLimit{Context: model.ContextLength, Output: model.ContextLength}
		}

		supportedInputModalities := filterLLMModalities(model.Architecture.InputModalities)
		supportedOutputModalities := filterLLMModalities(model.Architecture.OutputModalities)

		if len(supportedInputModalities) > 0 && len(supportedOutputModalities) > 0 {
			entry.Modalities = &OpencodeOutputModalities{
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

		models[entry.ID] = entry
	}

	for _, modelID := range providerConfig.Models {
		if matchedModelIDs[modelID] {
			continue
		}

		fmt.Fprintf(os.Stderr, "%s model %q was not found for provider %q, using ID as name\n", aurora.Yellow("warn:").String(), modelID, providerConfig.Name)
		models[modelID] = OpencodeOutputModel{ID: modelID, Name: modelID}
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
