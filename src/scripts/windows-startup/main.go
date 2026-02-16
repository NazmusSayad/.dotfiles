package main

import (
	helpers "dotfiles/src/helpers"
	"fmt"
)

type LaunchConfig struct {
	Path  string
	Args  []string
	Wait  bool
	Admin bool
}

func main() {
	launchConfigs := helpers.ReadConfig[[]LaunchConfig]("@/config/launch.jsonc")

	for _, config := range launchConfigs {
		resolvedCommand := helpers.ResolvePath(config.Path)
		fmt.Println("Starting: (", config.Admin, ")", resolvedCommand)

		helpers.ExecNativeCommand(
			append([]string{resolvedCommand}, config.Args...),
			helpers.ExecCommandOptions{
				AsAdmin: config.Admin,
				NoWait:  !config.Wait,
			},
		)
	}
}
