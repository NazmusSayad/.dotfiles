package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
	"gopkg.in/yaml.v3"
)

type ReadConfigOptions struct {
	Silent bool
}

func ReadConfig[T any](path string, options ...ReadConfigOptions) T {
	opts := ReadConfigOptions{}
	if len(options) > 0 {
		opts = options[0]
	} else if len(options) > 1 {
		panic("only one options struct is allowed")
	}

	resolvedPath := ResolvePath(path)
	if !opts.Silent {
		fmt.Println(aurora.Faint("CONFIG: " + resolvedPath))
	}

	f, err := os.Open(resolvedPath)
	if err != nil {
		panic("CONFIG: failed to open file")
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic("CONFIG: failed to read file")
	}

	if strings.HasSuffix(resolvedPath, ".json") || strings.HasSuffix(resolvedPath, ".jsonc") {
		var result T

		jsonBytes := jsonc.ToJSON(data)
		parserErr := json.Unmarshal(jsonBytes, &result)
		if parserErr != nil {
			panic("CONFIG JSON: failed to unmarshal")
		}

		return result
	}

	if strings.HasSuffix(resolvedPath, ".yaml") || strings.HasSuffix(resolvedPath, ".yml") {
		var result T

		parserErr := yaml.Unmarshal(data, &result)
		if parserErr != nil {
			panic("CONFIG YAML: failed to unmarshal")
		}

		return result
	}

	panic("CONFIG: invalid file type")
}
