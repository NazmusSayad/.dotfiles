package main

import (
	helpers "dotfiles/src/helpers"
	"fmt"
)

type LaunchConfig struct {
	Name  string
	Path  string
	Args  []string
	Skip  bool
	Admin bool
}

func main() {
	launchConfigs := helpers.ReadConfig[[]LaunchConfig]("@/config/launch.jsonc")

	for _, config := range launchConfigs {
		if config.Skip {
			fmt.Println("Skipping", config.Name)
			continue
		}

		resolvedCommand := helpers.ResolvePath(config.Path)
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
