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

func main() {
	TEMPLATE, err := os.ReadFile(filepath.Join(constants.SOURCE_DIR, "compile-alias", "template", "main.go"))
	if err != nil {
		panic(err)
	}

	if !helpers.IsFileExists(constants.BUILD_TEMP_DIR) {
		os.MkdirAll(constants.BUILD_TEMP_DIR, 0755)
	}

	for aliasName, aliasArguments := range constants.BIN_ALIASES {
		aliasCommand := aliasArguments[0]
		aliasArguments := aliasArguments[1:]

		TEMPLATE_CONTENT := strings.Replace(string(TEMPLATE), "{COMMAND}", aliasCommand, 1)
		TEMPLATE_CONTENT = strings.Replace(TEMPLATE_CONTENT, "{ARGUMENTS}", strings.Join(aliasArguments, "\", \""), 1)

		tempScriptPath := filepath.Join(constants.BUILD_TEMP_DIR, aliasName+".go")
		if err := os.WriteFile(tempScriptPath, []byte(TEMPLATE_CONTENT), 0644); err != nil {
			panic(err)
		}

		buildOutputPath := filepath.Join(constants.BUILD_SCRIPTS_DIR, aliasName+".exe")
		buildErr := helpers.ExecNativeCommand([]string{
			"go", "build", "-o", buildOutputPath, tempScriptPath,
		})

		if err := os.Remove(tempScriptPath); err != nil {
			panic(err)
		}

		if buildErr != nil {
			panic(buildErr)
		}

		fmt.Println(aurora.Faint("> Successfully compiled alias: ").String() + aliasName)
	}
}
