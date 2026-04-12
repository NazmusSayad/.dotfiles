package main

import (
	"encoding/json"
	"fmt"

	helpers "dotfiles/src/helpers"
	scoop "dotfiles/src/helpers/scoop"
)

func main() {
	msys2Apps := helpers.GetMsysApps()
	msys2JSON, err := json.MarshalIndent(msys2Apps, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msys2JSON))

	scoopApps := helpers.GetScoopApps()
	scoopJSON, err := json.MarshalIndent(scoopApps, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(scoopJSON))

	wingetApps := helpers.GetWingetApps()
	wingetJSON, err := json.MarshalIndent(wingetApps, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(wingetJSON))

	scoopOriginalApps := scoop.ReadScoopAppConfig()
	scoopOriginalJSON, err := json.MarshalIndent(scoopOriginalApps, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(scoopOriginalJSON))
}
