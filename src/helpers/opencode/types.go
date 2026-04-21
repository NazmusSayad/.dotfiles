package opencode

import "encoding/json"

type AuthProvider struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type AuthConfig map[string]AuthProvider

type OpencodeProviderConfig struct {
	Name    string `yaml:"name"`
	BaseURL string `yaml:"apiURL"`

	ModelsURL    string `yaml:"modelsURL"`
	ModelPrefix  string `yaml:"modelPrefix"`

	HasTurboMode bool `yaml:"hasTurboMode"`

	Models []string `yaml:"models"`
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
	ID         string                     `json:"id"`
	Name       string                     `json:"name"`
	Limit      *OpencodeOutputLimit       `json:"limit,omitempty"`
	Modalities *OpencodeOutputModalities  `json:"modalities,omitempty"`
	Family     string                     `json:"family,omitempty"`
	Variants   map[string]json.RawMessage `json:"variants,omitempty"`
}

type OpencodeOutputModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type OpencodeOutputLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

type OpencodeOutputProvider struct {
	API    string          `json:"api,omitempty"`
	Name   string          `json:"name,omitempty"`
	Models json.RawMessage `json:"models"`
}
