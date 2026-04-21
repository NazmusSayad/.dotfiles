package opencode

import "encoding/json"

type AuthProvider struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type AuthConfig map[string]AuthProvider

type OpencodeProviderConfig struct {
	ModelsURL string                        `yaml:"modelsURL"`
	Models    []OpencodeProviderConfigModel `yaml:"models"`
}

type OpencodeProviderConfigModel struct {
	ID         string `json:"id"`
	ContextCap int    `json:"context,omitempty"`
}

type OpenAiCompatibleModelsResponse struct {
	Data []OpenAiCompatibleModel `json:"data"`
}

type OpenAiCompatibleModel struct {
	ID            string                       `json:"id"`
	Name          string                       `json:"name"`
	ContextLength int                          `json:"context_length"`
	Architecture  OpenAiCompatibleArchitecture `json:"architecture"`
	Opencode      KiloOptionalOpencode         `json:"opencode"`
}

type OpenAiCompatibleArchitecture struct {
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
}

type KiloOptionalOpencode struct {
	Family   string                     `json:"family"`
	Variants map[string]json.RawMessage `json:"variants"`
}

type OpencodeOutputModel struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Family      string                     `json:"family,omitempty"`
	Attachment  *bool                      `json:"attachment,omitempty"`
	Reasoning   *bool                      `json:"reasoning,omitempty"`
	ToolCall    *bool                      `json:"tool_call,omitempty"`
	Temperature *bool                      `json:"temperature,omitempty"`
	ReleaseDate string                     `json:"release_date,omitempty"`
	Modalities  *OpencodeOutputModalities  `json:"modalities,omitempty"`
	Cost        *OpencodeOutputCost        `json:"cost,omitempty"`
	Limit       *OpencodeOutputLimit       `json:"limit,omitempty"`
	Variants    map[string]json.RawMessage `json:"variants,omitempty"`
}

type OpencodeOutputModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type OpencodeOutputLimit struct {
	Context int `json:"context,omitempty"`
	Input   int `json:"input,omitempty"`
	Output  int `json:"output,omitempty"`
}

type OpencodeOutputCost struct {
	Input     float64 `json:"input,omitempty"`
	Output    float64 `json:"output,omitempty"`
	CacheRead float64 `json:"cache_read,omitempty"`
}

type OpencodeOutputProvider struct {
	API    string          `json:"api,omitempty"`
	Name   string          `json:"name,omitempty"`
	Models json.RawMessage `json:"models"`
}
