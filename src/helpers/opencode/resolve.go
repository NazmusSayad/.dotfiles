package opencode

import (
	"fmt"

	"dotfiles/src/utils"

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

	if providerConfig.URL != "" {
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
		fetchedModel, hasModelInFetched := fetchedModels[configuredModel.ID]

		if hasModelInFetched && hasModelInModelsDotDev {
			resolvedModel := ApplyModelContextCap(modelsDevModel, configuredModel.ContextCap)

			if fetchedModel.Variants != nil {
				resolvedModel.Variants = fetchedModel.Variants
			}

			resolvedModelsMap[configuredModel.ID] = resolvedModel
			fmt.Println(aurora.Green("✓"), configuredModel.ID)
			continue
		}

		if hasModelInModelsDotDev {
			resolvedModelsMap[configuredModel.ID] = ApplyModelContextCap(modelsDevModel, configuredModel.ContextCap)
			fmt.Println(aurora.Faint("✓"), configuredModel.ID)
			continue
		}

		if hasModelInFetched {
			resolvedModelsMap[configuredModel.ID] = ApplyModelContextCap(fetchedModel, configuredModel.ContextCap)
			fmt.Println("✓", configuredModel.ID)
			continue
		}

		fmt.Println(aurora.Red("✗"), configuredModel.ID)
	}

	whitelist := make([]string, 0)
	for _, configuredModel := range providerConfig.Models {
		whitelist = append(whitelist, configuredModel.ID)
	}

	return OpencodeOutputProviderConfig{
		Models:    resolvedModelsMap,
		Whitelist: utils.SortArrayOfString(whitelist),
	}, nil
}
