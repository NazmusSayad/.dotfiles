package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

func main() {
	providerConfigs := helpers.ReadConfig[map[string]opencode.OpencodeProviderConfig]("@/config/ai/opencode-providers.yaml")
	authConfigPath := helpers.ResolvePath("~/.local/share/opencode/auth.json")
	authConfig := helpers.ReadConfig[opencode.AuthConfig](authConfigPath)

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
	outputProviderConfig := make(map[string]opencode.OpencodeOutputProviderConfig)

	fmt.Println()

	for providerID, providerConfig := range providerConfigs {
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

	if outputAgentModels.SmallModel != "" {
		fmt.Println(aurora.Green("Setting small model to:"), aurora.Yellow(outputAgentModels.SmallModel))
		fullConfig["small_model"] = outputAgentModels.SmallModel
	} else {
		fmt.Println(aurora.Faint("Unsetting small model"))
		delete(fullConfig, "small_model")
	}

	writeAgentModelConfig(fullConfig, "general", outputAgentModels.GeneralModel)
	writeAgentModelConfig(fullConfig, "explore", outputAgentModels.ExploreModel)
	writeAgentModelConfig(fullConfig, "compaction", outputAgentModels.CompactModel)

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

	fmt.Println(aurora.Green("Successfully updated OpenCode models!"))
}

func writeAgentModelConfig(fullConfig map[string]any, agent string, modelId string) {
	if modelId == "" {
		fmt.Println(aurora.Faint("Unsetting " + agent + " model"))
		deleteAgentModelConfig(fullConfig, agent)
		return
	}

	fmt.Println(aurora.Green("Setting "+agent+" model to:"), aurora.Yellow(modelId))

	if fullConfig["agent"] == nil {
		fullConfig["agent"] = make(map[string]any)
	}

	if fullConfig["agent"].(map[string]any)[agent] == nil {
		fullConfig["agent"].(map[string]any)[agent] = make(map[string]any)
	}

	fullConfig["agent"].(map[string]any)[agent].(map[string]any)["model"] = modelId
}

func deleteAgentModelConfig(fullConfig map[string]any, agent string) {
	if fullConfig["agent"] == nil {
		return
	}

	if fullConfig["agent"].(map[string]any)[agent] == nil {
		return
	}

	if fullConfig["agent"].(map[string]any)[agent].(map[string]any)["model"] == nil {
		return
	}

	delete(fullConfig["agent"].(map[string]any)[agent].(map[string]any), "model")
}
