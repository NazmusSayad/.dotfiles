package main

import (
	helpers "dotfiles/src"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LaunchConfig struct {
	Name  string
	Path  string
	Args  []string
	Skip  bool
	Admin bool
}

func main() {
	exePath, exePathErr := os.Executable()
	if exePathErr != nil {
		fmt.Println("Error getting executable path...")
		os.Exit(1)
	}

	cwd := filepath.Dir(exePath)
	fmt.Printf("CWD: %s\n", cwd)
	fullPath := filepath.Join(cwd, "./config/launch.jsonc")
	data, err := helpers.ReadJsoncAsJson(fullPath)
	if err != nil {
		fmt.Println("Error reading JSON file...")
		os.Exit(1)
	}

	var launchConfigs []LaunchConfig
	if err := json.Unmarshal(data, &launchConfigs); err != nil {
		fmt.Println("Error unmarshalling JSON into LaunchConfig...")
		os.Exit(1)
	}

	for _, config := range launchConfigs {
		if config.Skip {
			fmt.Println("Skipping", config.Name)
			continue
		}

		resolvedCommand := helpers.ResolvePath(cwd, config.Path)
		fmt.Println("Starting: (", config.Admin, ")", config.Name, resolvedCommand)

		if config.Admin {
			err := helpers.DetachedElevate(resolvedCommand, config.Args...)
			if err != nil {
				fmt.Println("Error elevating", config.Name)
				continue
			}
		} else {
			err := helpers.DetachedExec(resolvedCommand, config.Args...)
			if err != nil {
				fmt.Println("Error executing", config.Name)
				continue
			}
		}
	}
}
