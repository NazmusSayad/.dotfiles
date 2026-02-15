package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/logrusorgru/aurora/v4"
	"github.com/tidwall/jsonc"
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
		fmt.Println(aurora.Faint("JSON: " + resolvedPath))
	}

	f, err := os.Open(resolvedPath)
	if err != nil {
		panic("JSON: failed to open file")
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic("JSON: failed to read file")
	}

	var result T
	jsonBytes := jsonc.ToJSON(data)
	parserErr := json.Unmarshal(jsonBytes, &result)
	if parserErr != nil {
		panic("JSON: failed to unmarshal")
	}

	return result
}
