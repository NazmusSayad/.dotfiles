package main

import (
	"fmt"
	"os"
	"strings"

	"dotfiles/src/helpers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: alias <shell>\n")
		fmt.Fprintf(os.Stderr, "Supported shells: bash, pwsh\n")
		os.Exit(1)
	}

	shell := os.Args[1]
	aliases := helpers.ReadConfig[map[string]string]("@/config/alias.jsonc", helpers.ReadConfigOptions{Silent: true})

	var output strings.Builder

	switch shell {
	case "sh":
		for name, cmd := range aliases {
			output.WriteString("alias " + name + "=\"" + resolveUnixCMD(cmd) + "\"\n")
		}

	case "pwsh":
		for name, cmd := range aliases {
			output.WriteString("function " + name + " { " + resolveWindowsCMD(cmd) + " @args }\n")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown shell: %s\n", shell)
		os.Exit(1)
	}

	fmt.Print(output.String())
}

func resolveUnixCMD(cmd string) string {
	return os.ExpandEnv(cmd)
}

func resolveWindowsCMD(cmd string) string {
	return os.ExpandEnv(cmd)
}
