package opencode

import (
	"strings"

	"dotfiles/src/helpers"
)

type providersConfig = map[string]OpencodeProviderConfig

const settingsPrefix = "~"

func ReadConfig() (providersConfig, OpencodeSettingsConfig) {
	providerConfigs := helpers.ReadConfig[providersConfig]("@/config/ai/opencode-providers.yaml")
	opencodeConfig := helpers.ReadConfig[OpencodeSettingsConfig]("@/config/ai/opencode-providers.yaml")

	for providerId := range providerConfigs {
		if strings.HasPrefix(providerId, settingsPrefix) {
			delete(providerConfigs, providerId)
		}
	}

	return providerConfigs, opencodeConfig
}
