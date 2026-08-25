package symlink

import (
	"fmt"
	"runtime"
	"sort"

	helpers "dotfiles/src/helpers"
	"gopkg.in/yaml.v3"
)

type StringOrArray []string

func (s *StringOrArray) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var target string
		if err := value.Decode(&target); err != nil {
			return err
		}
		*s = []string{target}
		return nil
	case yaml.SequenceNode:
		return value.Decode((*[]string)(s))
	default:
		return fmt.Errorf("value must be a string or array of strings: %s", value.Value)
	}
}

type Config struct {
	Source      string
	LinkTargets []string
	CopyTargets []string
}

type Entry struct {
	Link    StringOrArray `yaml:"-"`
	Win     StringOrArray `yaml:"Win"`
	Mac     StringOrArray `yaml:"Mac"`
	Copy    StringOrArray `yaml:"Copy"`
	WinCopy StringOrArray `yaml:"Win.Copy"`
	MacCopy StringOrArray `yaml:"Mac.Copy"`
}

func (e *Entry) UnmarshalYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode, yaml.SequenceNode:
		return node.Decode(&e.Link)
	case yaml.MappingNode:
		type plainEntry Entry
		return node.Decode((*plainEntry)(e))
	default:
		return fmt.Errorf("invalid symlink entry: %s", node.Value)
	}
}

func (e Entry) resolve() ([]string, []string) {
	linkTargets := e.Link
	copyTargets := e.Copy

	switch runtime.GOOS {
	case "windows":
		if len(e.Win) > 0 {
			linkTargets = e.Win
		}
		if len(e.WinCopy) > 0 {
			copyTargets = e.WinCopy
		}
	case "darwin":
		if len(e.Mac) > 0 {
			linkTargets = e.Mac
		}
		if len(e.MacCopy) > 0 {
			copyTargets = e.MacCopy
		}
	}

	return []string(linkTargets), []string(copyTargets)
}

func ReadConfigs() []Config {
	rawConfigs := helpers.ReadConfig[map[string]Entry]("@/config/symlink.yml")

	sources := make([]string, 0, len(rawConfigs))
	for source := range rawConfigs {
		sources = append(sources, source)
	}
	sort.Strings(sources)

	var configs []Config
	for _, source := range sources {
		linkTargets, copyTargets := rawConfigs[source].resolve()
		if len(linkTargets) == 0 && len(copyTargets) == 0 {
			continue
		}

		configs = append(configs, Config{
			Source:      source,
			LinkTargets: linkTargets,
			CopyTargets: copyTargets,
		})
	}

	return configs
}
