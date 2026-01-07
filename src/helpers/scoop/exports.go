package scoop

import (
	"encoding/json"
	"os/exec"
	"time"
)

type Bucket struct {
	Name      string
	Source    string
	Updated   time.Time
	Manifests int
}

type App struct {
	Updated time.Time
	Source  string
	Info    string
	Version string
	Name    string
}

type Export struct {
	Buckets []Bucket `json:"buckets"`
	Apps    []App    `json:"apps"`
}

func GetScoopExports() (Export, error) {
	cmd := exec.Command("scoop", "export")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Export{}, err
	}

	var export Export
	err = json.Unmarshal(output, &export)
	return export, err
}

func GetScoopApps() ([]App, error) {
	exports, err := GetScoopExports()
	if err != nil {
		return nil, err
	}
	return exports.Apps, nil
}

func GetScoopBuckets() ([]Bucket, error) {
	exports, err := GetScoopExports()
	if err != nil {
		return nil, err
	}
	return exports.Buckets, nil
}
