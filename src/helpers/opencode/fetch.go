package opencode

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"slices"
	"strings"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

type modelsDotDevProvider struct {
	Models map[string]OpencodeStandardModel `json:"models"`
}

var allowedModalities = []string{"text", "audio", "image", "video", "pdf"}

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

func FetchModels(providerID string, providerURL string, auth *AuthProvider) (map[string]OpencodeStandardModel, error) {
	fmt.Printf("%s %s\n", aurora.Yellow("Fetching models from").String(), aurora.Faint(providerURL).String())

	req, err := http.NewRequest("GET", providerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s models: %w", providerID, err)
	}

	if auth != nil && auth.Type == "api" && auth.Key != "" {
		req.Header.Set("Authorization", "Bearer "+auth.Key)
		fmt.Printf("%s Using API key from auth config\n", aurora.Faint("Using auth:").String())
	} else {
		envApiKey := os.Getenv(strings.ToUpper(providerID) + "_API_KEY")
		if envApiKey != "" {
			req.Header.Set("Authorization", "Bearer "+envApiKey)
			fmt.Printf("%s Using %s environment variable\n", aurora.Faint("Using auth:"), aurora.Faint(strings.ToUpper(providerID)+"_API_KEY").String())
		}
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

	return models, nil
}

func FetchOpenrouterModels(auth AuthConfig) (map[string]OpencodeStandardModel, error) {
	openrouterAuth := auth["openrouter"]
	if openrouterAuth.Type == "api" && openrouterAuth.Key != "" {
		return FetchModels("openrouter", OPENROUTER_MODELS_URL, &openrouterAuth)
	}

	return FetchModels("openrouter", OPENROUTER_MODELS_URL, nil)
}

func FetchModelsDotDev() (map[string]map[string]OpencodeStandardModel, error) {
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

	var payload map[string]modelsDotDevProvider
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode models.dev response: %w", err)
	}

	result := make(map[string]map[string]OpencodeStandardModel)
	for providerID, provider := range payload {
		providerModels := make(map[string]OpencodeStandardModel)
		maps.Copy(providerModels, provider.Models)

		result[providerID] = providerModels
	}

	return result, nil
}
