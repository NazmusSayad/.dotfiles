package opencode

import (
	"fmt"
	"os"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

type OpencodeResolveAgentModels struct {
	SmallModel   string
	GeneralModel string
	ExploreModel string
	CompactModel string
}

func ResolveOpencodeProvider(
	providerId string, providerConfig OpencodeProviderConfig, modelsDotDevProvider ModelsDotDevProvider,
	openrouterModels map[string]OpencodeStandardModel, currentAgentModels OpencodeResolveAgentModels, authConfig AuthConfig,
) (OpencodeOutputProviderConfig, OpencodeResolveAgentModels, error) {
	var fetchedModels map[string]OpencodeStandardModel

	if providerConfig.URL != "" {
		apiKey := ResolveApiKey(providerId, modelsDotDevProvider, authConfig)
		if models, err := FetchModels(providerId, providerConfig.URL, apiKey); err == nil {
			fetchedModels = models
		} else {
			fmt.Printf("%s Failed to fetch models for %s: %s\n", aurora.Red("Error:"), providerId, err.Error())
			return OpencodeOutputProviderConfig{}, OpencodeResolveAgentModels{}, err
		}
	}

	resolvedModelsMap := make(map[string]OpencodeStandardModel)
	for _, configuredModel := range providerConfig.Models {
		currentAgentModels = resolveAgentModel(providerId, configuredModel, currentAgentModels)

		if providerConfig.WhitelistOnly {
			continue
		}

		openrouterModel, hasModelInOpenrouter := openrouterModels[configuredModel.OpenrouterModelId]
		modelsDevModel, hasModelInModelsDotDev := modelsDotDevProvider.Models[configuredModel.ID]
		fetchedModel, hasModelInFetched := fetchedModels[configuredModel.ID]

		var resolvedModel *OpencodeStandardModel

		if hasModelInModelsDotDev && (hasModelInFetched || hasModelInOpenrouter) {
			resolvedModel = &modelsDevModel

			if fetchedModel.Variants != nil {
				resolvedModel.Variants = fetchedModel.Variants
			}

			fmt.Println(aurora.Green("[ALL]"), configuredModel.ID)

		} else if hasModelInModelsDotDev {
			resolvedModel = &modelsDevModel
			fmt.Println(aurora.Green("[MDD]"), configuredModel.ID)

		} else if hasModelInOpenrouter {
			resolvedModel = &openrouterModel
			resolvedModel.ID = configuredModel.ID
			fmt.Println(aurora.Blue("[OPR]"), configuredModel.ID)

		} else if hasModelInFetched {
			resolvedModel = &fetchedModel
			fmt.Println(aurora.Cyan("[API]"), configuredModel.ID)
		}

		if resolvedModel != nil {
			if configuredModel.Nitro {
				resolvedModel.ID = utils.Ternary(configuredModel.Nitro, resolvedModel.ID+":nitro", resolvedModel.ID)
			}

			if len(configuredModel.Headers) > 0 {
				if resolvedModel.Headers == nil {
					resolvedModel.Headers = make(map[string]string)
				}
				for k, v := range configuredModel.Headers {
					resolvedModel.Headers[k] = v
				}
			}

			resolvedModelsMap[configuredModel.ID] = applyModelContextCap(*resolvedModel, configuredModel.ContextCap)
		} else {
			fmt.Println(aurora.Red("[ERR]"), configuredModel.ID)
			resolvedModelsMap[configuredModel.ID] = OpencodeStandardModel{
				ID:      utils.Ternary(configuredModel.Nitro, configuredModel.ID+":nitro", configuredModel.ID),
				Name:    configuredModel.ID,
				Headers: configuredModel.Headers,
			}
		}
	}

	whitelist := make([]string, 0)
	for _, configuredModel := range providerConfig.Models {
		whitelist = append(whitelist, configuredModel.ID)
	}

	if providerConfig.WhitelistOnly {
		fmt.Println(aurora.Faint("Only whitelisted models will be included for this provider"))
		return OpencodeOutputProviderConfig{
			Whitelist: utils.SortArrayOfString(whitelist),
		}, currentAgentModels, nil
	}

	return OpencodeOutputProviderConfig{
		Models:    resolvedModelsMap,
		Whitelist: utils.SortArrayOfString(whitelist),
	}, currentAgentModels, nil
}

func resolveAgentModel(providerId string, modelConfig OpencodeProviderConfigModel, currentAgentModels OpencodeResolveAgentModels) OpencodeResolveAgentModels {
	if modelConfig.AsSmallModel {
		if currentAgentModels.SmallModel != "" {
			fmt.Printf(
				"%s Multiple models marked as small model. Models %s and %s will be used as the small model.\n",
				aurora.Red("ERROR:"), currentAgentModels.SmallModel, modelConfig.ID,
			)
			os.Exit(1)
		}

		currentAgentModels.SmallModel = providerId + "/" + modelConfig.ID
	}

	if modelConfig.AsGeneralModel {
		if currentAgentModels.GeneralModel != "" {
			fmt.Printf(
				"%s Multiple models marked as general model. Models %s and %s will be used as the general model.\n",
				aurora.Red("ERROR:"), currentAgentModels.GeneralModel, modelConfig.ID,
			)
			os.Exit(1)
		}

		currentAgentModels.GeneralModel = providerId + "/" + modelConfig.ID
	}

	if modelConfig.AsExploreModel {
		if currentAgentModels.ExploreModel != "" {
			fmt.Printf(
				"%s Multiple models marked as explore model. Models %s and %s will be used as the explore model.\n",
				aurora.Red("ERROR:"), currentAgentModels.ExploreModel, modelConfig.ID,
			)
			os.Exit(1)
		}

		currentAgentModels.ExploreModel = providerId + "/" + modelConfig.ID
	}

	if modelConfig.AsCompactModel {
		if currentAgentModels.CompactModel != "" {
			fmt.Printf(
				"%s Multiple models marked as compact model. Models %s and %s will be used as the compact model.\n",
				aurora.Red("ERROR:"), currentAgentModels.CompactModel, modelConfig.ID,
			)
			os.Exit(1)
		}

		currentAgentModels.CompactModel = providerId + "/" + modelConfig.ID
	}

	return currentAgentModels
}

func applyModelContextCap(model OpencodeStandardModel, contextCap int) OpencodeStandardModel {
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
