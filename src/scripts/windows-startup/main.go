package main

import (
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
)

type LaunchConfig struct {
	Path string
	Args []string
	Wait bool

	AsUser  bool
	AsAdmin bool
}

func main() {
	launchConfigs := helpers.ReadConfig[[]LaunchConfig]("@/config/launch.jsonc")

	for _, config := range launchConfigs {
		resolvedCommand := helpers.ResolvePath(config.Path)
		fmt.Println("Starting: (", config.AsAdmin, ")", resolvedCommand)

		resolvedArguments := make([]string, len(config.Args))
		for i, arg := range config.Args {
			resolvedArguments[i] = os.ExpandEnv(arg)
		}

		helpers.ExecNativeCommand(
			append([]string{resolvedCommand}, resolvedArguments...),
			helpers.ExecCommandOptions{
				NoWait:      !config.Wait,
				AsAdmin:     config.AsAdmin,
				AsGsudoUser: config.AsUser,
			},
		)
	}
}
