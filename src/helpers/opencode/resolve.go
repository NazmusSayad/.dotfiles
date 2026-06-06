package opencode

import (
	"fmt"
	"os"

	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
)

type OpencodeResolveAgentModels struct {
	MainModel   string
	SmallModel  string
	AgentsModel map[string]string
}

func ResolveOpencodeProvider(
	providerId string, providerConfig OpencodeProviderConfig, modelsDotDevProvider ModelsDotDevProvider,
	openrouterModels map[string]OpencodeStandardModel, currentAgentModels OpencodeResolveAgentModels, authConfig AuthConfig,
) (OpencodeStandardProvider, OpencodeResolveAgentModels, error) {
	var fetchedModels map[string]OpencodeStandardModel

	if providerConfig.URL != "" {
		apiKey := ResolveApiKey(providerId, modelsDotDevProvider, authConfig)
		if models, err := FetchModels(providerId, providerConfig.URL, apiKey); err == nil {
			fetchedModels = models
		} else {
			fmt.Printf("%s Failed to fetch models for %s: %s\n", aurora.Red("Error:"), providerId, err.Error())
			return OpencodeStandardProvider{}, OpencodeResolveAgentModels{}, err
		}
	}

	resolvedModelsMap := make(map[string]OpencodeStandardModel)
	for _, modelConfig := range providerConfig.Models {
		currentAgentModels = resolveAgentModel(providerId, modelConfig, currentAgentModels)

		openrouterModel, hasModelInOpenrouter := openrouterModels[modelConfig.OpenrouterModelId]
		modelsDevModel, hasModelInModelsDotDev := modelsDotDevProvider.Models[modelConfig.ID]
		fetchedModel, hasModelInFetched := fetchedModels[modelConfig.ID]

		if hasModelInModelsDotDev && !providerConfig.Include && !modelConfig.Include &&
			modelConfig.Variants == nil && modelConfig.Options == nil && len(modelConfig.Headers) == 0 &&
			!modelConfig.Nitro && modelConfig.ContextCap == 0 {
			fmt.Println(aurora.Green("[MDD]"), modelConfig.ID)
			continue
		}

		var resolvedModel *OpencodeStandardModel

		if hasModelInModelsDotDev && (hasModelInFetched || hasModelInOpenrouter) {
			resolvedModel = &modelsDevModel

			if fetchedModel.Variants != nil {
				resolvedModel.Variants = fetchedModel.Variants
			}

			fmt.Println(aurora.Green("[ALL]"), modelConfig.ID)

		} else if hasModelInModelsDotDev {
			resolvedModel = &modelsDevModel
			fmt.Println(aurora.Green("[MDD]"), modelConfig.ID)

		} else if hasModelInOpenrouter {
			resolvedModel = &openrouterModel
			resolvedModel.ID = modelConfig.ID
			fmt.Println(aurora.Blue("[OPR]"), modelConfig.ID)

		} else if hasModelInFetched {
			resolvedModel = &fetchedModel
			fmt.Println(aurora.Cyan("[API]"), modelConfig.ID)
		}

		if resolvedModel == nil {
			fmt.Println(aurora.Red("[ERR]"), modelConfig.ID)
			resolvedModel = &OpencodeStandardModel{
				ID:      modelConfig.ID,
				Name:    modelConfig.ID,
				Headers: modelConfig.Headers,
			}
		}

		if modelConfig.Nitro {
			resolvedModel.ID = utils.Ternary(modelConfig.Nitro, resolvedModel.ID+":nitro", resolvedModel.ID)
		}

		if modelConfig.Options != nil {
			resolvedModel.Options = modelConfig.Options
		}

		if modelConfig.Variants != nil {
			resolvedModel.Variants = modelConfig.Variants
		}

		if len(modelConfig.Headers) > 0 {
			if resolvedModel.Headers == nil {
				resolvedModel.Headers = make(map[string]string)
			}
			for k, v := range modelConfig.Headers {
				resolvedModel.Headers[k] = v
			}
		}

		resolvedModelsMap[modelConfig.ID] = applyModelContextCap(*resolvedModel, modelConfig.ContextCap)
	}

	whitelist := make([]string, 0)
	for _, configuredModel := range providerConfig.Models {
		whitelist = append(whitelist, configuredModel.ID)
	}

	return OpencodeStandardProvider{
		Models:    utils.Ternary(len(resolvedModelsMap) > 0, resolvedModelsMap, nil),
		Whitelist: utils.SortArrayOfString(whitelist),
	}, currentAgentModels, nil
}

func resolveAgentModel(providerId string, modelConfig OpencodeProviderConfigModel, currentAgentModels OpencodeResolveAgentModels) OpencodeResolveAgentModels {
	modelId := providerId + "/" + modelConfig.ID

	if modelConfig.AsMain {
		if currentAgentModels.MainModel != "" {
			fmt.Printf(
				"%s Multiple models marked as agent model. Models %s and %s will be used as the agent model.\n",
				aurora.Red("ERROR:"), currentAgentModels.MainModel, modelConfig.ID,
			)
			os.Exit(1)
		}

		currentAgentModels.MainModel = modelId
	}

	if modelConfig.AsSmall {
		if currentAgentModels.SmallModel != "" {
			fmt.Printf(
				"%s Multiple models marked as small model. Models %s and %s will be used as the small model.\n",
				aurora.Red("ERROR:"), currentAgentModels.SmallModel, modelConfig.ID,
			)
			os.Exit(1)
		}

		currentAgentModels.SmallModel = modelId
	}

	if modelConfig.AsAgentTitle {
		currentAgentModels = setAgentModel(currentAgentModels, "title", modelId)
	}

	if modelConfig.AsAgentGeneral {
		currentAgentModels = setAgentModel(currentAgentModels, "general", modelId)
	}

	if modelConfig.AsAgentExplore {
		currentAgentModels = setAgentModel(currentAgentModels, "explore", modelId)
	}

	if modelConfig.AsAgentSummary {
		currentAgentModels = setAgentModel(currentAgentModels, "summary", modelId)
	}

	if modelConfig.AsAgentCompaction {
		currentAgentModels = setAgentModel(currentAgentModels, "compaction", modelId)
	}

	return currentAgentModels
}

func setAgentModel(currentAgentModels OpencodeResolveAgentModels, agentId string, modelId string) OpencodeResolveAgentModels {
	if currentAgentModels.AgentsModel == nil {
		currentAgentModels.AgentsModel = make(map[string]string)
	}

	if currentAgentModels.AgentsModel[agentId] != "" {
		fmt.Printf(
			"%s Multiple models marked as %s model. Models %s and %s will be used as the %s model.\n",
			aurora.Red("ERROR:"), agentId, currentAgentModels.AgentsModel[agentId], modelId, agentId,
		)
		os.Exit(1)
	}

	currentAgentModels.AgentsModel[agentId] = modelId
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
