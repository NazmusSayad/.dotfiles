package linuxpkg

import (
	"dotfiles/src/helpers"

	"gopkg.in/yaml.v3"
)

type snapApp struct {
	ID      string `yaml:"id"`
	Channel string `yaml:"channel"`
	Classic bool   `yaml:"classic"`
}

func (app *snapApp) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		app.ID = node.Value
		return nil
	}

	type plainApp snapApp
	return node.Decode((*plainApp)(app))
}

type debApp struct {
	ID     string `yaml:"id"`
	URL    string `yaml:"url"`
	GitHub string `yaml:"github"`
	Asset  string `yaml:"asset"`
}

type appsConfig struct {
	Apt  []string  `yaml:"apt"`
	Snap []snapApp `yaml:"snap"`
	Deb  []debApp  `yaml:"deb"`
}

func readConfig() appsConfig {
	return helpers.ReadConfig[appsConfig]("@/config/apps.linux.yml")
}
