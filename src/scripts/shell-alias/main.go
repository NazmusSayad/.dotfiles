package main

import (
	"fmt"
	"os"
	"os/exec"
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
			output.WriteString("alias " + name + "=\"" + cygpath(cmd, "-u") + "\"\n")
		}

	case "pwsh":
		for name, cmd := range aliases {
			output.WriteString("function " + name + " { " + cygpath(cmd, "-w") + " @args }\n")
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown shell: %s\n", shell)
		os.Exit(1)
	}

	fmt.Print(output.String())
}

func cygpath(input string, mode string) string {
	expanded := os.ExpandEnv(input)
	if expanded == input {
		return input
	}

	cmd := exec.Command("cygpath", mode, expanded)
	out, err := cmd.Output()
	if err != nil {
		return expanded
	}

	return strings.TrimSpace(string(out))
}
