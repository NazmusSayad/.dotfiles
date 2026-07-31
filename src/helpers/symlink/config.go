package symlink

import (
	"encoding/json"
	"fmt"
	"runtime"

	helpers "dotfiles/src/helpers"
)

type StringOrArray []string

func (s *StringOrArray) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}

	var array []string
	if err := json.Unmarshal(data, &array); err == nil {
		*s = array
		return nil
	}

	return fmt.Errorf("value must be a string or array of strings: %s", string(data))
}

type Config struct {
	Copy    bool
	Source  string
	Targets []string
}

type rawConfig struct {
	Copy      bool          `json:"Copy"`
	Source    string        `json:"Source"`
	Target    StringOrArray `json:"Target"`
	TargetWin StringOrArray `json:"Target.Win"`
	TargetMac StringOrArray `json:"Target.Mac"`
}

func ReadConfigs() []Config {
	rawConfigs := helpers.ReadConfig[[]rawConfig]("@/config/symlink.jsonc")

	var configs []Config
	for _, raw := range rawConfigs {
		targets := resolveTargets(raw)
		if len(targets) == 0 {
			continue
		}

		configs = append(configs, Config{
			Copy:    raw.Copy,
			Source:  raw.Source,
			Targets: targets,
		})
	}

	return configs
}

func resolveTargets(raw rawConfig) []string {
	if runtime.GOOS == "windows" && len(raw.TargetWin) > 0 {
		return raw.TargetWin
	}

	if runtime.GOOS == "darwin" && len(raw.TargetMac) > 0 {
		return raw.TargetMac
	}

	if len(raw.Target) > 0 {
		return raw.Target
	}

	return nil
}
