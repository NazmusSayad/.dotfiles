package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

func main() {
	providersConfig, opencodeConfig := opencode.ReadConfig()
	authConfigPath := helpers.ResolvePath("~/.local/share/opencode/auth.json")
	authConfig := helpers.ReadConfig[opencode.AuthConfig](authConfigPath, helpers.ReadConfigOptions{SkipError: true})

	modelsDotDevResponse, modelsDotDevError := opencode.FetchModelsDotDev()
	if modelsDotDevError != nil {
		fmt.Println("failed to fetch models.dev models:", modelsDotDevError)
		return
	}

	openrouterModelsResponse, openrouterModelsError := opencode.FetchOpenrouterModels(authConfig)
	if openrouterModelsError != nil {
		fmt.Println("failed to fetch openrouter models:", openrouterModelsError)
		return
	}

	outputAgentModels := opencode.OpencodeResolveAgentModels{}
	outputProviderConfig := make(map[string]opencode.OpencodeStandardProvider)

	fmt.Println()

	for providerID, providerConfig := range providersConfig {
		fmt.Printf("%s %s\n", aurora.Blue("Syncing models for"), aurora.Bold(providerID))

		result, resolvedAgentModels, err := opencode.ResolveOpencodeProvider(
			providerID, providerConfig,
			modelsDotDevResponse[providerID],
			openrouterModelsResponse,
			outputAgentModels,
			authConfig,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s failed to resolve provider %q: %v\n", aurora.Yellow("warn:"), providerID, err)
			fmt.Println()
			continue
		}

		if providerConfig.API != "" {
			result.API = providerConfig.API
		}

		if len(providerConfig.Env) > 0 {
			result.Env = providerConfig.Env
		}

		outputAgentModels = resolvedAgentModels
		outputProviderConfig[providerID] = result
		fmt.Println()
	}

	enabledProviders := make([]string, 0)
	for providerID := range outputProviderConfig {
		enabledProviders = append(enabledProviders, providerID)
	}

	configPath := helpers.ResolvePath("@/config/ai/opencode.json")
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("failed to read opencode config:", err)
		os.Exit(1)
	}

	var fullConfig map[string]any
	if err := json.Unmarshal(jsonc.ToJSON(configBytes), &fullConfig); err != nil {
		fmt.Println("failed to decode opencode config:", err)
		os.Exit(1)
	}

	fullConfig["provider"] = outputProviderConfig
	fullConfig["enabled_providers"] = utils.SortArrayOfString(enabledProviders)

	if outputAgentModels.MainModel != "" {
		fmt.Println(aurora.Green("Setting main model to:"), aurora.Yellow(outputAgentModels.MainModel))
		fullConfig["model"] = outputAgentModels.MainModel
	} else {
		fmt.Println(aurora.Faint("Unsetting main model"))
		delete(fullConfig, "model")
	}

	if outputAgentModels.SmallModel != "" {
		fmt.Println(aurora.Green("Setting small model to:"), aurora.Yellow(outputAgentModels.SmallModel))
		fullConfig["small_model"] = outputAgentModels.SmallModel
	} else {
		fmt.Println(aurora.Faint("Unsetting small model"))
		delete(fullConfig, "small_model")
	}

	fullConfig["agent"] = opencodeConfig.Agents
	setAgentModel(fullConfig, "title", outputAgentModels.AgentsModel["title"], outputAgentModels.AgentsOptions["title"])
	setAgentModel(fullConfig, "general", outputAgentModels.AgentsModel["general"], outputAgentModels.AgentsOptions["general"])
	setAgentModel(fullConfig, "explore", outputAgentModels.AgentsModel["explore"], outputAgentModels.AgentsOptions["explore"])
	setAgentModel(fullConfig, "summary", outputAgentModels.AgentsModel["summary"], outputAgentModels.AgentsOptions["summary"])
	setAgentModel(fullConfig, "compaction", outputAgentModels.AgentsModel["compaction"], outputAgentModels.AgentsOptions["compaction"])

	newConfigBytes, err := json.Marshal(fullConfig)
	if err != nil {
		fmt.Println("failed to encode config:", err)
		os.Exit(1)
	}

	mergedConfigRaw, err := helpers.MergeJSONObject(string(configBytes), string(newConfigBytes))
	if err != nil {
		fmt.Println("failed to merge config:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(configPath, []byte(mergedConfigRaw), 0o644); err != nil {
		fmt.Println("failed to write opencode config:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(aurora.Cyan("Refreshing opencode models..."))
	refreshErr := helpers.ExecNativeCommand(
		[]string{"opencode", "models", "--refresh"},
		helpers.ExecCommandOptions{Silent: true},
	)
	if refreshErr != nil {
		fmt.Println("failed to refresh opencode models")
	}

	fmt.Println(aurora.Cyan("Formatting opencode config..."))
	prettierErr := helpers.ExecNativeCommand(
		[]string{"npx", "-y", "prettier", "--write", configPath},
		helpers.ExecCommandOptions{Silent: true},
	)
	if prettierErr != nil {
		fmt.Println("failed to format opencode config")
	}

	fmt.Println(aurora.Green("Successfully updated OpenCode config!"))
}

func setAgentModel(fullConfig map[string]any, agent string, modelId string, options map[string]any) {
	prevConfig := fullConfig["agent"].(map[string]any)[agent]
	resolvedConfig := make(map[string]any)

	if modelId != "" {
		fmt.Println(aurora.Green("Setting "+agent+" model to:"), aurora.Yellow(modelId))
		resolvedConfig["model"] = modelId
	}

	if prevConfig != nil {
		maps.Copy(resolvedConfig, prevConfig.(map[string]any))
	}

	if len(options) > 0 {
		agentOptions, _ := resolvedConfig["options"].(map[string]any)
		resolvedConfig["options"] = mergeAgentOptions(options, agentOptions)
	}

	if len(resolvedConfig) > 0 {
		fullConfig["agent"].(map[string]any)[agent] = resolvedConfig
	} else {
		delete(fullConfig["agent"].(map[string]any), agent)
	}
}

func mergeAgentOptions(modelOptions, agentOptions map[string]any) map[string]any {
	merged := maps.Clone(modelOptions)
	if merged == nil {
		merged = make(map[string]any)
	}

	for key, value := range agentOptions {
		modelNested, modelIsMap := merged[key].(map[string]any)
		agentNested, agentIsMap := value.(map[string]any)
		if modelIsMap && agentIsMap {
			merged[key] = mergeAgentOptions(modelNested, agentNested)
		} else {
			merged[key] = value
		}
	}

	return merged
}
