package main

import (
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/utils"
)

type LaunchConfig struct {
	Path string
	Args []string

	Dir  string
	Wait bool

	AsUser  bool
	AsAdmin bool
}

func main() {
	launchConfigs := helpers.ReadConfig[[]LaunchConfig]("@/config/launch.jsonc")

	for _, config := range launchConfigs {
		resolvedCommand := helpers.ResolvePath(config.Path)
		fmt.Println("Starting: (", config.AsAdmin, ")", resolvedCommand)

		resolvedWorkingDir := utils.Ternary(config.Dir != "", helpers.ResolvePath(config.Dir), "")
		fmt.Println("Working dir:", resolvedWorkingDir)

		resolvedArguments := make([]string, len(config.Args))
		for i, arg := range config.Args {
			resolvedArguments[i] = os.ExpandEnv(arg)
		}

		helpers.ExecNativeCommand(
			append([]string{resolvedCommand}, resolvedArguments...),
			helpers.ExecCommandOptions{
				NoWait:      !config.Wait,
				Dir:         resolvedWorkingDir,
				AsAdmin:     config.AsAdmin,
				AsGsudoUser: config.AsUser,
			},
		)
	}
}
