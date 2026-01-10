package scoop

import (
	"dotfiles/src/utils"
	"encoding/json"
	"os/exec"
)

type ScoopBucket struct {
	Name string
}

type ScoopApp struct {
	Name    string
	Source  string
	Version string
}

type ScoopExport struct {
	Buckets []ScoopBucket `json:"buckets"`
	Apps    []ScoopApp    `json:"apps"`
}

func GetScoopExports() ScoopExport {
	cmd := exec.Command("scoop", "export")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic("Error getting scoop exports: " + err.Error())
	}

	var export ScoopExport
	err = json.Unmarshal(output, &export)
	if err != nil {
		panic("Error unmarshalling scoop exports: " + err.Error())
	}

	return export
}

func GetScoopExportAppMap(export ScoopExport) map[string]ScoopApp {
	appMap := make(map[string]ScoopApp)

	for _, app := range export.Apps {
		appMap[app.Source+"/"+app.Name] = app
	}

	return appMap
}

func GetScoopExportBucketsList(export ScoopExport) []string {
	bucketList := []string{}

	for _, bucket := range export.Buckets {
		bucketList = append(bucketList, bucket.Name)
	}

	return utils.UniqueArray(bucketList)
}
