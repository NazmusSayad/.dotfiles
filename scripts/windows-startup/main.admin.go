package main

import (
	helpers "dotfiles/src/helpers"
	"encoding/json"
	"os"
)

type LaunchConfig struct {
	Name  string
	Path  string
	Args  []string
	Skip  bool
	Admin bool
}

func main() {
	data, err := helpers.ReadDotfilesConfigJSONC("./config/launch.jsonc")
	if err != nil {
		println("Error reading JSON file...")
		os.Exit(1)
	}

	var launchConfigs []LaunchConfig
	if err := json.Unmarshal(data, &launchConfigs); err != nil {
		println("Error unmarshalling JSON into LaunchConfig...")
		os.Exit(1)
	}

	for _, config := range launchConfigs {
		if config.Skip {
			println("Skipping", config.Name)
			continue
		}

		resolvedCommand := helpers.ResolvePath(config.Path)
		println("Starting: (", config.Admin, ")", config.Name, resolvedCommand)

		if config.Admin {
			err := helpers.DetachedElevate(resolvedCommand, config.Args...)
			if err != nil {
				println("Error elevating", config.Name)
				continue
			}
		} else {
			err := helpers.DetachedExec(resolvedCommand, config.Args...)
			if err != nil {
				println("Error executing", config.Name)
				continue
			}
		}
	}
}
