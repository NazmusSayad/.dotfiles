package main

import (
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
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

		resolvedArguments := make([]string, len(config.Args))
		for i, arg := range config.Args {
			resolvedArguments[i] = os.ExpandEnv(arg)
		}

		helpers.ExecNativeCommand(
			append([]string{resolvedCommand}, resolvedArguments...),
			helpers.ExecCommandOptions{
				AsAdmin: config.Admin,
				NoWait:  !config.Wait,
			},
		)
	}
}
