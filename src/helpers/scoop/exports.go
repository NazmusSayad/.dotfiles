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

func GetScoopExports() Export {
	cmd := exec.Command("scoop", "export")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic("Error getting scoop exports: " + err.Error())
	}

	var export Export
	err = json.Unmarshal(output, &export)
	if err != nil {
		panic("Error unmarshalling scoop exports: " + err.Error())
	}

	return export
}
