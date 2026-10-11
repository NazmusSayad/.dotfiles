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
	InheritPerm bool
	Create      string
}

type Entry struct {
	Link        StringOrArray `yaml:"Link"`
	Copy        StringOrArray `yaml:"Copy"`
	LinkWin     StringOrArray `yaml:"LinkWin"`
	LinkMac     StringOrArray `yaml:"LinkMac"`
	LinkWsl     StringOrArray `yaml:"LinkWsl"`
	CopyWin     StringOrArray `yaml:"CopyWin"`
	CopyMac     StringOrArray `yaml:"CopyMac"`
	CopyWsl     StringOrArray `yaml:"CopyWsl"`
	Create      string        `yaml:"Create"`
	InheritPerm bool          `yaml:"InheritPerm"`
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
		if len(e.LinkWin) > 0 {
			linkTargets = e.LinkWin
		}
		if len(e.CopyWin) > 0 {
			copyTargets = e.CopyWin
		}
	case "darwin":
		if len(e.LinkMac) > 0 {
			linkTargets = e.LinkMac
		}
		if len(e.CopyMac) > 0 {
			copyTargets = e.CopyMac
		}
	case "linux":
		if helpers.IsWSL() {
			if len(e.LinkWsl) > 0 {
				linkTargets = e.LinkWsl
			}
			if len(e.CopyWsl) > 0 {
				copyTargets = e.CopyWsl
			}
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

		switch rawConfigs[source].Create {
		case "", "Folder", "File":
		default:
			panic(fmt.Sprintf("SYMLINK: invalid Create value %q for %s", rawConfigs[source].Create, source))
		}

		configs = append(configs, Config{
			Source:      source,
			LinkTargets: linkTargets,
			CopyTargets: copyTargets,
			InheritPerm: rawConfigs[source].InheritPerm,
			Create:      rawConfigs[source].Create,
		})
	}

	return configs
}
