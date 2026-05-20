package opencode

import (
	"fmt"
	"reflect"
	"strings"

	"dotfiles/src/helpers"

	"github.com/logrusorgru/aurora/v4"
)

type providersConfig = map[string]OpencodeProviderConfig

const settingsPrefix = "~"

func ReadConfig() (providersConfig, OpencodeSettingsConfig) {
	providerConfigs := helpers.ReadConfig[providersConfig]("@/config/ai/opencode-providers.yaml")

	opencodeConfig := OpencodeSettingsConfig{}
	v := reflect.ValueOf(&opencodeConfig).Elem()

	for providerId, config := range providerConfigs {
		if strings.HasPrefix(providerId, settingsPrefix) {
			delete(providerConfigs, providerId)

			normalizedKey := strings.TrimPrefix(providerId, settingsPrefix)
			field := v.FieldByName(normalizedKey)

			if field.IsValid() && field.CanSet() {
				field.Set(reflect.ValueOf(config))
			} else {
				fmt.Println(aurora.Red("error:"), aurora.Yellow(fmt.Sprintf("invalid config key %q in opencode config", providerId)))
			}
		}
	}

	return providerConfigs, opencodeConfig
}
