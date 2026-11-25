package linuxpkg

import (
	"slices"

	"dotfiles/src/helpers"

	"gopkg.in/yaml.v3"
)

type aptApp struct {
	ID         string `yaml:"id"`
	Recommends bool   `yaml:"recommends"`
}

func (app *aptApp) UnmarshalYAML(node *yaml.Node) error {
	app.Recommends = true
	if node.Kind == yaml.ScalarNode {
		app.ID = node.Value
		return nil
	}

	type plainApp aptApp
	return node.Decode((*plainApp)(app))
}

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
	Apt        []aptApp  `yaml:"apt"`
	AptExclude []string  `yaml:"apt-exclude"`
	Snap       []snapApp `yaml:"snap"`
	Deb        []debApp  `yaml:"deb"`
}

func readConfig() appsConfig {
	config := helpers.ReadConfig[appsConfig]("@/config/apps.linux.yml")

	for _, id := range config.AptExclude {
		if slices.Contains(config.aptIDs(), id) {
			panic("CONFIG: " + id + " is in both apt and apt-exclude")
		}
	}

	return config
}

func (config appsConfig) aptIDs() []string {
	ids := []string{}
	for _, app := range config.Apt {
		ids = append(ids, app.ID)
	}

	return ids
}
