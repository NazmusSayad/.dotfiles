package opencode

import (
	"fmt"

	"github.com/logrusorgru/aurora/v4"
)

func ApplyModelContextCap(model OpencodeStandardModel, contextCap int) OpencodeStandardModel {
	if model.Cost == nil {
		model.Cost = &OpencodeStandardCost{
			Input:     0,
			Output:    0,
			CacheRead: 0,
		}
	}

	if model.Cost.Input == 0 || model.Cost.Output == 0 {
		model.Cost = nil
	}

	if contextCap <= 0 || model.Limit == nil {
		return model
	}

	if model.Limit.Context > contextCap {
		model.Limit.Context = contextCap
	}

	if model.Limit.Input > contextCap {
		model.Limit.Input = contextCap
	}

	if model.Limit.Output > contextCap {
		model.Limit.Output = contextCap
	}

	return model
}

func ResolveOpencodeProvider(providerId string, providerConfig OpencodeProviderConfig, modelsDotDevProvider map[string]OpencodeStandardModel, authConfig AuthConfig) (OpencodeOutputProviderConfig, error) {
	var fetchedModels map[string]OpencodeStandardModel

	if providerConfig.ModelsURL != "" {
		var providerAuth *AuthProvider
		if auth, ok := authConfig[providerId]; ok {
			providerAuth = &auth
		}

		models, err := FetchModels(providerId, providerConfig, providerAuth)
		if err != nil {
			return OpencodeOutputProviderConfig{}, err
		}

		fetchedModels = models
	}

	resolvedModelsMap := make(map[string]OpencodeStandardModel)
	for _, configuredModel := range providerConfig.Models {
		modelsDevModel, hasModelInModelsDotDev := modelsDotDevProvider[configuredModel.ID]
		if hasModelInModelsDotDev {
			resolvedModelsMap[configuredModel.ID] = ApplyModelContextCap(modelsDevModel, configuredModel.ContextCap)
			fmt.Println(aurora.Faint("✓"), configuredModel.ID)
			continue
		}

		fetchedModel, hasModelInFetched := fetchedModels[configuredModel.ID]
		if hasModelInFetched {
			resolvedModelsMap[configuredModel.ID] = ApplyModelContextCap(fetchedModel, configuredModel.ContextCap)
			fmt.Println(aurora.Green("✓"), configuredModel.ID)
			continue
		}

		fmt.Println(aurora.Red("✗"), configuredModel.ID)
	}

	whitelist := make([]string, 0, len(providerConfig.Models))
	for _, configuredModel := range providerConfig.Models {
		whitelist = append(whitelist, configuredModel.ID)
	}

	return OpencodeOutputProviderConfig{Models: resolvedModelsMap, Whitelist: whitelist}, nil
}
