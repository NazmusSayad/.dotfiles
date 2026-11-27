package linuxpkg

import (
	"bytes"
	"errors"
	"os"

	"dotfiles/src/helpers"

	"gopkg.in/yaml.v3"
)

const lockPath = "@/.local/apps.linux.lock"

type lockEntry struct {
	Version      string   `yaml:"version,omitempty"`
	Channel      string   `yaml:"channel,omitempty"`
	Classic      bool     `yaml:"classic,omitempty"`
	NoRecommends bool     `yaml:"noRecommends,omitempty"`
	URL          string   `yaml:"url,omitempty"`
	GitHub       string   `yaml:"github,omitempty"`
	AptSources   []string `yaml:"aptSources,omitempty"`
	Uninstall    string   `yaml:"uninstall,omitempty"`
	InstalledAt  string   `yaml:"installedAt"`
}

type appsLock struct {
	Apt    map[string]lockEntry `yaml:"apt"`
	Snap   map[string]lockEntry `yaml:"snap"`
	Deb    map[string]lockEntry `yaml:"deb"`
	Script map[string]lockEntry `yaml:"script"`
}

func readLock() appsLock {
	data, err := os.ReadFile(helpers.ResolvePath(lockPath))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic("LOCK: failed to read file: " + err.Error())
	}

	lock := appsLock{}
	if err := yaml.Unmarshal(data, &lock); err != nil {
		panic("LOCK: failed to unmarshal: " + err.Error())
	}

	if lock.Apt == nil {
		lock.Apt = map[string]lockEntry{}
	}
	if lock.Snap == nil {
		lock.Snap = map[string]lockEntry{}
	}
	if lock.Deb == nil {
		lock.Deb = map[string]lockEntry{}
	}
	if lock.Script == nil {
		lock.Script = map[string]lockEntry{}
	}

	return lock
}

func writeLock(lock appsLock) {
	var data bytes.Buffer
	encoder := yaml.NewEncoder(&data)
	encoder.SetIndent(2)
	if err := encoder.Encode(lock); err != nil {
		panic("LOCK: failed to marshal: " + err.Error())
	}
	if err := encoder.Close(); err != nil {
		panic("LOCK: failed to marshal: " + err.Error())
	}

	if err := os.WriteFile(helpers.ResolvePath(lockPath), data.Bytes(), 0o644); err != nil {
		panic("LOCK: failed to write file: " + err.Error())
	}
}
