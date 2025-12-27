package main

import (
	"dotfiles/src/constants"
	"dotfiles/src/helpers"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora/v4"
)

type StructName struct {
	Name      string
	Arguments []string
}

var ALIASES = []StructName{
	{
		Name:      "fsc",
		Arguments: []string{"fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive"},
	},

	{
		Name:      "ghp",
		Arguments: []string{"gh", "pr", "create", "-B"},
	},

	{
		Name:      "ghv",
		Arguments: []string{"gh", "repo", "view"},
	},

	{
		Name:      "ghw",
		Arguments: []string{"gh", "repo", "view", "--web"},
	},

	{
		Name:      "gds",
		Arguments: []string{"git", "diff", "--stat"},
	},

	{
		Name:      "gp",
		Arguments: []string{"git", "pull"},
	},
}

func main() {
	TEMPLATE, err := os.ReadFile(filepath.Join(constants.SOURCE_DIR, "compile-alias", "template", "main.go"))
	if err != nil {
		panic(err)
	}

	if !helpers.IsFileExists(constants.BUILD_TEMP_DIR) {
		os.MkdirAll(constants.BUILD_TEMP_DIR, 0755)
	}

	for _, alias := range ALIASES {
		aliasCommand := alias.Arguments[0]
		aliasArguments := alias.Arguments[1:]

		TEMPLATE_CONTENT := strings.Replace(string(TEMPLATE), "{COMMAND}", aliasCommand, 1)
		TEMPLATE_CONTENT = strings.Replace(TEMPLATE_CONTENT, "{ARGUMENTS}", strings.Join(aliasArguments, "\", \""), 1)

		tempScriptPath := filepath.Join(constants.BUILD_TEMP_DIR, alias.Name+".go")
		if err := os.WriteFile(tempScriptPath, []byte(TEMPLATE_CONTENT), 0644); err != nil {
			panic(err)
		}

		buildOutputPath := filepath.Join(constants.BUILD_SCRIPTS_DIR, alias.Name+".exe")
		buildErr := helpers.ExecNativeCommand(helpers.ExecCommandOptions{
			Command: "go",
			Args:    []string{"build", "-o", buildOutputPath, tempScriptPath},
		})

		if err := os.Remove(tempScriptPath); err != nil {
			panic(err)
		}

		if buildErr != nil {
			panic(buildErr)
		}

		fmt.Println(aurora.Faint("> Successfully compiled alias: ").String() + alias.Name)
	}
}
