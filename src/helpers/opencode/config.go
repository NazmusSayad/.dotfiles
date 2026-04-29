package opencode

import (
	"fmt"
	"strings"

	"dotfiles/src/helpers"
)

type rawOpencodeProviderConfig struct {
	URL    string `yaml:"url"`
	Models []any  `yaml:"models"`
}

const CONTEXT_CAP = 400000

func ReadOpencodeProvidersConfig() map[string]OpencodeProviderConfig {
	rawProviders := helpers.ReadConfig[map[string]rawOpencodeProviderConfig]("@/config/ai/opencode-providers.yaml")
	providers := make(map[string]OpencodeProviderConfig, len(rawProviders))

	for providerID, rawProvider := range rawProviders {
		models := make([]OpencodeProviderConfigModel, 0, len(rawProvider.Models))

		for i, rawModel := range rawProvider.Models {
			switch model := rawModel.(type) {
			case string:
				modelID := strings.TrimSpace(model)
				if modelID == "" {
					panic(fmt.Sprintf("opencode providers config: provider %q model at index %d is empty", providerID, i))
				}

				models = append(models, OpencodeProviderConfigModel{ID: modelID, ContextCap: CONTEXT_CAP})
			case map[string]any:
				modelIDValue, ok := model["id"]
				if !ok {
					panic(fmt.Sprintf("opencode providers config: provider %q model at index %d is missing id", providerID, i))
				}

				modelID, ok := modelIDValue.(string)
				if !ok {
					panic(fmt.Sprintf("opencode providers config: provider %q model at index %d has invalid id", providerID, i))
				}

				modelID = strings.TrimSpace(modelID)
				if modelID == "" {
					panic(fmt.Sprintf("opencode providers config: provider %q model at index %d is empty", providerID, i))
				}

				contextCap := CONTEXT_CAP
				if contextValue, exists := model["context"]; exists {
					switch context := contextValue.(type) {
					case int:
						contextCap = context
					case int64:
						contextCap = int(context)
					case float64:
						if context != float64(int(context)) {
							panic(fmt.Sprintf("opencode providers config: provider %q model %q has non-integer context", providerID, modelID))
						}
						contextCap = int(context)
					default:
						panic(fmt.Sprintf("opencode providers config: provider %q model %q has invalid context", providerID, modelID))
					}

					if contextCap <= 0 {
						panic(fmt.Sprintf("opencode providers config: provider %q model %q has invalid context", providerID, modelID))
					}
				}

				openrouterId := ""
				if openrouterIdValue, exists := model["openrouterId"]; exists {
					if routerId, ok := openrouterIdValue.(string); ok {
						openrouterId = strings.TrimSpace(routerId)
					} else {
						panic(fmt.Sprintf("opencode providers config: provider %q model %q has invalid openrouterId", providerID, modelID))
					}
				}

				nitro := false
				if nitroValue, exists := model["nitro"]; exists {
					if n, ok := nitroValue.(bool); ok {
						nitro = n
					} else {
						panic(fmt.Sprintf("opencode providers config: provider %q model %q has invalid nitro", providerID, modelID))
					}
				}

				models = append(models, OpencodeProviderConfigModel{ID: modelID, ContextCap: contextCap, OpenrouterModelId: openrouterId, Nitro: nitro})
			default:
				panic(fmt.Sprintf("opencode providers config: provider %q model at index %d must be a string or object", providerID, i))
			}
		}

		providers[providerID] = OpencodeProviderConfig{
			URL: rawProvider.URL, Models: models,
		}
	}

	return providers
}
