package opencode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

var (
	allowedModalities  = []string{"text", "audio", "image", "video", "pdf"}
	fetchedModelsCache = make(map[string]map[string]OpencodeStandardModel)
)

const OPENROUTER_MODELS_URL = "https://openrouter.ai/api/v1/models"

func filterLLMModalities(modalities []string) []string {
	var filtered []string
	for _, m := range modalities {
		if slices.Contains(allowedModalities, m) {
			filtered = append(filtered, m)
		}
	}

	return filtered
}

func FetchModels(providerID string, providerURL string, apiKey string) (map[string]OpencodeStandardModel, error) {
	if cached, ok := fetchedModelsCache[providerURL]; ok {
		fmt.Println(aurora.Yellow("Getting cached models from:"), aurora.Faint(providerURL))
		return cached, nil
	}

	fmt.Println(
		aurora.Yellow("Fetching models from:"), aurora.Faint(providerURL),
		utils.Ternary(apiKey != "", aurora.Green("Authenticated").String(), aurora.Red("Unauthenticated").String()),
	)

	req, err := http.NewRequest("GET", providerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s models: %w", providerID, err)
	}

	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s models: %w", providerID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s models: %s", providerID, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s response body: %w", providerID, err)
	}

	var payload OpenAiCompatibleModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", providerID, err)
	}

	models := map[string]OpencodeStandardModel{}

	for _, model := range payload.Data {
		entry := OpencodeStandardModel{
			ID:   model.ID,
			Name: utils.Ternary(model.Name == "", model.ID, model.Name),
		}

		if model.ContextLength > 0 {
			entry.Limit = &OpencodeStandardLimit{
				Context: model.ContextLength,
				Input:   model.ContextLength,
				Output:  model.ContextLength,
			}
		}

		supportedInputModalities := filterLLMModalities(model.Architecture.InputModalities)
		supportedOutputModalities := filterLLMModalities(model.Architecture.OutputModalities)

		if len(supportedInputModalities) > 0 && len(supportedOutputModalities) > 0 {
			entry.Modalities = &OpencodeStandardModalities{
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

	fetchedModelsCache[providerURL] = models
	return models, nil
}

func FetchOpenrouterModels(auth AuthConfig) (map[string]OpencodeStandardModel, error) {
	return FetchModels(
		"openrouter", OPENROUTER_MODELS_URL,
		ResolveApiKey("openrouter", ModelsDotDevProvider{Env: []string{"OPENROUTER_API_KEY"}}, auth),
	)
}

func FetchModelsDotDev() (map[string]ModelsDotDevProvider, error) {
	resp, err := http.Get("https://models.dev/api.json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models.dev API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch models.dev API: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read models.dev response body: %w", err)
	}

	var payload map[string]ModelsDotDevProvider
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode models.dev response: %w", err)
	}

	return payload, nil
}
