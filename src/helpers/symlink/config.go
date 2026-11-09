package symlink

import (
	"fmt"
	"maps"
	"os"
	"runtime"
	"slices"
	"strings"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/utils"
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
	Expanded    bool
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
	entries := map[string][]Entry{}
	for _, path := range []string{"@/config/symlink.yml", "@/config/symlink.claude.yml", "#/symlink.yml", "#/symlink.claude.yml"} {
		if strings.HasPrefix(path, "#/") && !utils.IsFileExists(helpers.ResolvePath(path)) {
			continue
		}
		for source, entry := range helpers.ReadConfig[map[string]Entry](path) {
			entries[source] = append(entries[source], entry)
		}
	}

	var configs []Config
	for _, source := range slices.Sorted(maps.Keys(entries)) {
		config := Config{Source: source, Create: entries[source][0].Create, InheritPerm: entries[source][0].InheritPerm}
		for _, entry := range entries[source] {
			if entry.Create != config.Create || entry.InheritPerm != config.InheritPerm {
				panic(fmt.Sprintf("SYMLINK: conflicting Create or InheritPerm for %s", source))
			}
			linkTargets, copyTargets := entry.resolve()
			config.LinkTargets = append(config.LinkTargets, linkTargets...)
			config.CopyTargets = append(config.CopyTargets, copyTargets...)
		}
		if len(config.LinkTargets) == 0 && len(config.CopyTargets) == 0 {
			continue
		}

		switch config.Create {
		case "", "Folder", "File":
		default:
			panic(fmt.Sprintf("SYMLINK: invalid Create value %q for %s", config.Create, source))
		}

		if strings.HasSuffix(source, "/*") {
			configs = append(configs, expandFolder(config)...)
		} else {
			configs = append(configs, config)
		}
	}

	targetSources := map[string]string{}
	for _, config := range configs {
		for _, target := range slices.Concat(config.LinkTargets, config.CopyTargets) {
			if other, exists := targetSources[target]; exists {
				panic(fmt.Sprintf("SYMLINK: %s is a target of both %s and %s", target, other, config.Source))
			}
			targetSources[target] = config.Source
		}
	}

	return configs
}

func expandFolder(config Config) []Config {
	if config.Create != "" {
		panic(fmt.Sprintf("SYMLINK: Create is not allowed with %s", config.Source))
	}

	folder := strings.TrimSuffix(config.Source, "*")
	items, err := os.ReadDir(helpers.ResolvePath(folder))
	if err != nil {
		panic(fmt.Sprintf("SYMLINK: cannot read %s: %v", folder, err))
	}

	var configs []Config
	for _, item := range items {
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}

		itemConfig := Config{Source: folder + item.Name(), InheritPerm: config.InheritPerm, Expanded: true}
		for _, target := range config.LinkTargets {
			itemConfig.LinkTargets = append(itemConfig.LinkTargets, target+"/"+item.Name())
		}
		for _, target := range config.CopyTargets {
			itemConfig.CopyTargets = append(itemConfig.CopyTargets, target+"/"+item.Name())
		}
		configs = append(configs, itemConfig)
	}

	return configs
}
