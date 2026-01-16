package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

type MiseEnv struct {
	Source string `json:"source"`
	Tool   string `json:"tool"`
	Value  string `json:"value"`
}

func main() {
	miseEnv := map[string]MiseEnv{}
	miseEnvCmd := exec.Command("mise", "env", "--json-extended")
	output, err := miseEnvCmd.Output()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(output, &miseEnv)
	if err != nil {
		panic(err)
	}

	for name, env := range miseEnv {
		if name == "PATH" {
			paths := strings.Split(env.Value, ";")
			os.Stdout.WriteString(strings.Join(paths, "\n"))
			return
		}
	}
}
