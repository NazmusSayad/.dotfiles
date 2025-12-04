package main

import (
	helpers "dotfiles/src"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
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
			err := helpers.Elevate(resolvedCommand, config.Args...)
			if err != nil {
				fmt.Println("Error elevating", config.Name, err)
				continue
			}
		} else {
			cmd := exec.Command(resolvedCommand, config.Args...)
			cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008}
			cmd.Start()
		}
	}
}
