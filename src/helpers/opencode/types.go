package opencode

type AuthProvider struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type AuthConfig map[string]AuthProvider

type OpencodeSettingsConfig struct {
	Agents map[string]any `yaml:"~agents,omitempty"`
}

type OpencodeProviderConfig struct {
	URL string   `yaml:"url"`
	API string   `yaml:"api,omitempty"`
	Env []string `yaml:"env,omitempty"`

	Include bool                          `yaml:"include,omitempty"`
	Models  []OpencodeProviderConfigModel `yaml:"models"`
}

type OpencodeProviderConfigModel struct {
	ID                string `yaml:"id"`
	Nitro             bool   `yaml:"nitro,omitempty"`
	ContextCap        int    `yaml:"context,omitempty"`
	OpenrouterModelId string `yaml:"openrouterId,omitempty"`

	Include  bool              `yaml:"include,omitempty"`
	Options  any               `yaml:"options,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	Variants map[string]any    `yaml:"variants,omitempty"`

	AsMain  bool `yaml:"main,omitempty"`
	AsSmall bool `yaml:"small,omitempty"`

	AsAgentTitle      bool `yaml:"title,omitempty"`
	AsAgentGeneral    bool `yaml:"general,omitempty"`
	AsAgentExplore    bool `yaml:"explore,omitempty"`
	AsAgentSummary    bool `yaml:"summary,omitempty"`
	AsAgentCompaction bool `yaml:"compaction,omitempty"`
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
	Family   string         `json:"family"`
	Variants map[string]any `json:"variants"`
}

type OpencodeStandardModel struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Family      string                      `json:"family,omitempty"`
	Attachment  *bool                       `json:"attachment,omitempty"`
	Reasoning   *bool                       `json:"reasoning,omitempty"`
	ToolCall    *bool                       `json:"tool_call,omitempty"`
	Temperature *bool                       `json:"temperature,omitempty"`
	ReleaseDate string                      `json:"release_date,omitempty"`
	Modalities  *OpencodeStandardModalities `json:"modalities,omitempty"`
	Cost        *OpencodeStandardCost       `json:"cost,omitempty"`
	Limit       *OpencodeStandardLimit      `json:"limit,omitempty"`
	Options     any                         `json:"options,omitempty"`
	Headers     map[string]string           `json:"headers,omitempty"`
	Variants    map[string]any              `json:"variants,omitempty"`
}

type OpencodeStandardModalities struct {
	Input  []string `json:"input"`
	Output []string `json:"output"`
}

type OpencodeStandardLimit struct {
	Context int `json:"context,omitempty"`
	Input   int `json:"input,omitempty"`
	Output  int `json:"output,omitempty"`
}

type OpencodeStandardCost struct {
	Input     float64 `json:"input,omitempty"`
	Output    float64 `json:"output,omitempty"`
	CacheRead float64 `json:"cache_read,omitempty"`
}

type OpencodeStandardProvider struct {
	API       string                           `json:"api,omitempty"`
	Env       []string                         `json:"env,omitempty"`
	Models    map[string]OpencodeStandardModel `json:"models"`
	Whitelist []string                         `json:"whitelist"`
}

type ModelsDotDevProvider struct {
	ID     string                           `json:"id"`
	Name   string                           `json:"name"`
	Env    []string                         `json:"env"`
	ApiUrl string                           `json:"api"`
	Models map[string]OpencodeStandardModel `json:"models"`
}
