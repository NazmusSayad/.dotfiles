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
	Wait  bool
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

		helpers.ExecNativeCommand(
			append([]string{resolvedCommand}, config.Args...),
			helpers.ExecCommandOptions{
				AsAdmin: config.Admin == true,
				AsUser:  config.Admin == false,
				NoWait:  !config.Wait,
			},
		)
	}
}
