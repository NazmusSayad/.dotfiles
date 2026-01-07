package scoop

import "dotfiles/src/helpers"

type ScoopConfig struct {
	ID            string
	Name          string
	Bucket        string
	SkipHashCheck bool
}

func ReadScoopConfig() []ScoopConfig {
	return helpers.ReadConfig[[]ScoopConfig]("@/config/scoop-apps.jsonc")
}
