package scoop

import "dotfiles/src/helpers"

type ScoopAppConfig struct {
	ID            string
	Name          string
	Bucket        string
	Version       string
	SkipHashCheck bool
}

func ReadScoopAppConfig() []ScoopAppConfig {
	return helpers.ReadConfig[[]ScoopAppConfig]("@/config/scoop-apps.jsonc")
}
