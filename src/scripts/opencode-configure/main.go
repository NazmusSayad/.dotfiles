package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"dotfiles/src/helpers"
	"dotfiles/src/helpers/opencode"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
)

func main() {
	providersConfig, opencodeConfig := opencode.ReadConfig()
	modelsDotDevResponse, modelsDotDevError := opencode.FetchModelsDotDev()
	if modelsDotDevError != nil {
		fmt.Println("failed to fetch models.dev models:", modelsDotDevError)
		return
	}

	openrouterModelsResponse, openrouterModelsError := opencode.FetchOpenrouterModels()
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
		if len(providersConfig[providerID].Models) > 0 {
			enabledProviders = append(enabledProviders, providerID)
		}
	}

	configPath := helpers.ResolvePath("@/config/ai/opencode/opencode.json")
	compiledConfigPath := helpers.ResolvePath("@/config/ai/opencode/opencode.compile.json")
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

	if shell, hasShell := fullConfig["shell"]; hasShell {
		shellName, ok := shell.(string)
		if !ok {
			fmt.Println("invalid shell in opencode config:", shell)
			os.Exit(1)
		}
		shellPath, err := exec.LookPath(shellName)
		if err != nil {
			fmt.Println("failed to find shell:", err)
			os.Exit(1)
		}
		fmt.Println(aurora.Green("Setting shell to:"), aurora.Yellow(shellPath))
		fullConfig["shell"] = shellPath
	}

	providers := make(map[string]any)
	for providerID, provider := range outputProviderConfig {
		resolved := resolveProvider(provider)
		configuredModels := providersConfig[providerID].Models
		if len(configuredModels) > 0 {
			models, _ := resolved["models"].(map[string]any)
			if models == nil {
				models = make(map[string]any)
			}
			for catalogModelID, catalogModel := range modelsDotDevResponse[providerID].Models {
				catalogModelIDs := []string{catalogModelID}
				if catalogModel.Experimental != nil {
					for mode := range catalogModel.Experimental.Modes {
						catalogModelIDs = append(catalogModelIDs, catalogModelID+"-"+mode)
					}
				}
				for _, modelID := range catalogModelIDs {
					isConfigured := slices.ContainsFunc(configuredModels, func(model opencode.OpencodeProviderConfigModel) bool {
						return model.ID == modelID
					})
					if !isConfigured {
						models[modelID] = map[string]any{"disabled": true}
					}
				}
			}
			if len(models) > 0 {
				resolved["models"] = models
			}
		}
		providers[providerID] = resolved
	}
	fullConfig["providers"] = providers
	policies := []map[string]string{{"action": "provider.use", "resource": "*", "effect": "deny"}}
	for _, providerID := range utils.SortArrayOfString(enabledProviders) {
		policies = append(policies, map[string]string{"action": "provider.use", "resource": providerID, "effect": "allow"})
	}
	experimental, _ := fullConfig["experimental"].(map[string]any)
	if experimental == nil {
		experimental = make(map[string]any)
	}
	experimental["policies"] = policies
	fullConfig["experimental"] = experimental

	if outputAgentModels.MainModel != "" {
		fmt.Println(aurora.Green("Setting main model to:"), aurora.Yellow(outputAgentModels.MainModel))
		fullConfig["model"] = outputAgentModels.MainModel
	} else {
		fmt.Println(aurora.Faint("Unsetting main model"))
		delete(fullConfig, "model")
	}

	fullConfig["agents"] = opencodeConfig.Agents
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

	if err := os.WriteFile(compiledConfigPath, []byte(mergedConfigRaw), 0o644); err != nil {
		fmt.Println("failed to write compiled opencode config:", err)
		os.Exit(1)
	}

	cliConfigPath := helpers.ResolvePath("@/config/ai/opencode/cli.json")
	compiledCliConfigPath := helpers.ResolvePath("@/config/ai/opencode/cli.compile.json")
	cliConfigBytes, err := os.ReadFile(cliConfigPath)
	if err != nil {
		fmt.Println("failed to read opencode cli config:", err)
		os.Exit(1)
	}

	var cliConfig map[string]any
	if err := json.Unmarshal(jsonc.ToJSON(cliConfigBytes), &cliConfig); err != nil {
		fmt.Println("failed to decode opencode cli config:", err)
		os.Exit(1)
	}

	attention, _ := cliConfig["attention"].(map[string]any)
	sounds, _ := attention["sounds"].(map[string]any)
	for soundID, soundPath := range sounds {
		path, ok := soundPath.(string)
		if !ok {
			fmt.Printf("invalid sound path for %q: %v\n", soundID, soundPath)
			os.Exit(1)
		}
		if !filepath.IsAbs(path) {
			sounds[soundID] = filepath.ToSlash(filepath.Join(helpers.ResolvePath("~/.config/opencode"), path))
		}
	}

	newCliConfigBytes, err := json.Marshal(cliConfig)
	if err != nil {
		fmt.Println("failed to encode cli config:", err)
		os.Exit(1)
	}

	mergedCliConfigRaw, err := helpers.MergeJSONObject(string(cliConfigBytes), string(newCliConfigBytes))
	if err != nil {
		fmt.Println("failed to merge cli config:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(compiledCliConfigPath, []byte(mergedCliConfigRaw), 0o644); err != nil {
		fmt.Println("failed to write compiled opencode cli config:", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println(aurora.Cyan("Formatting compiled opencode configs..."))
	prettierErr := helpers.ExecNativeCommand(
		[]string{"mise", "exec", "--", "prettier", "--write", compiledConfigPath, compiledCliConfigPath},
		helpers.ExecCommandOptions{Silent: true},
	)
	if prettierErr != nil {
		fmt.Println("failed to format compiled opencode configs")
	}

	fmt.Println(aurora.Green("Successfully updated OpenCode config!"))
}

func setAgentModel(fullConfig map[string]any, agent string, modelId string, options map[string]any) {
	prevConfig := fullConfig["agents"].(map[string]any)[agent]
	resolvedConfig := make(map[string]any)

	if modelId != "" {
		fmt.Println(aurora.Green("Setting "+agent+" model to:"), aurora.Yellow(modelId))
		resolvedConfig["model"] = modelId
	}

	if prevConfig != nil {
		maps.Copy(resolvedConfig, prevConfig.(map[string]any))
	}

	if len(options) > 0 {
		request := make(map[string]any)
		if prevRequest, ok := resolvedConfig["request"].(map[string]any); ok {
			maps.Copy(request, prevRequest)
		}
		agentBody, _ := request["body"].(map[string]any)
		request["body"] = mergeAgentOptions(options, agentBody)
		resolvedConfig["request"] = request
	}

	if len(resolvedConfig) > 0 {
		fullConfig["agents"].(map[string]any)[agent] = resolvedConfig
	} else {
		delete(fullConfig["agents"].(map[string]any), agent)
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

func resolveProvider(provider opencode.OpencodeStandardProvider) map[string]any {
	result := make(map[string]any)
	if provider.API != "" {
		result["settings"] = map[string]any{"baseURL": provider.API}
	}
	if len(provider.Env) > 0 {
		result["env"] = provider.Env
	}
	if len(provider.Models) > 0 {
		models := make(map[string]any)
		for id, model := range provider.Models {
			entry := map[string]any{"name": model.Name}
			if model.ID != "" {
				entry["modelID"] = model.ID
			}
			if model.Family != "" {
				entry["family"] = model.Family
			}
			if model.Limit != nil {
				entry["limit"] = model.Limit
			}
			if model.Cost != nil {
				cost := map[string]any{"input": model.Cost.Input, "output": model.Cost.Output}
				if model.Cost.CacheRead > 0 {
					cost["cache"] = map[string]any{"read": model.Cost.CacheRead}
				}
				entry["cost"] = cost
			}
			if model.Modalities != nil || model.ToolCall != nil {
				capabilities := map[string]any{"tools": true, "input": []string{"text", "image"}, "output": []string{"text"}}
				if model.ToolCall != nil {
					capabilities["tools"] = *model.ToolCall
				}
				if model.Modalities != nil {
					capabilities["input"] = model.Modalities.Input
					capabilities["output"] = model.Modalities.Output
				}
				entry["capabilities"] = capabilities
			}
			if model.Options != nil {
				entry["settings"] = model.Options
			}
			if len(model.Headers) > 0 {
				entry["headers"] = model.Headers
			}
			if len(model.Variants) > 0 {
				variants := make([]map[string]any, 0, len(model.Variants))
				ids := make([]string, 0, len(model.Variants))
				for variantID := range model.Variants {
					ids = append(ids, variantID)
				}
				slices.Sort(ids)
				for _, variantID := range ids {
					variants = append(variants, map[string]any{"id": variantID, "settings": model.Variants[variantID]})
				}
				entry["variants"] = variants
			}
			models[id] = entry
		}
		result["models"] = models
	}
	return result
}
