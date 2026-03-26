package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"dotfiles/src/helpers"
)

type kiloModelsResponse struct {
	Data []kiloModel `json:"data"`
}

type kiloModel struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Opencode      kiloOpencode `json:"opencode"`
	ContextLength int          `json:"context_length"`
}

type kiloOpencode struct {
	Family   string                     `json:"family"`
	Variants map[string]json.RawMessage `json:"variants"`
}

type opencodeModel struct {
	ID       string                     `json:"id"`
	Name     string                     `json:"name"`
	Limit    *opencodeModelLimit        `json:"limit,omitempty"`
	Family   string                     `json:"family,omitempty"`
	Variants map[string]json.RawMessage `json:"variants,omitempty"`
}

type opencodeModelLimit struct {
	Context int `json:"context"`
	Output  int `json:"output"`
}

type opencodeConfig struct {
	Provider struct {
		KiloCustom struct {
			Whitelist []string `json:"whitelist"`
		} `json:"kilo-custom"`
	} `json:"provider"`
}

func main() {
	resp, err := http.Get("https://api.kilo.ai/api/gateway/models")
	if err != nil {
		fmt.Println("failed to fetch models:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("failed to fetch models:", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to read response body:", err)
		os.Exit(1)
	}

	var payload kiloModelsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		fmt.Println("failed to decode response:", err)
		os.Exit(1)
	}

	slices.SortFunc(payload.Data, func(a, b kiloModel) int {
		if a.Name < b.Name {
			return -1
		}

		if a.Name > b.Name {
			return 1
		}

		return 0
	})

	config := helpers.ReadConfig[opencodeConfig]("@/config/ai/opencode.json", helpers.ReadConfigOptions{Silent: true})
	whitelist := map[string]struct{}{}

	for _, id := range config.Provider.KiloCustom.Whitelist {
		whitelist[id] = struct{}{}
	}

	models := map[string]opencodeModel{}

	for _, model := range payload.Data {
		if _, ok := whitelist[model.ID]; !ok {
			continue
		}

		entry := opencodeModel{
			ID:   model.ID,
			Name: model.Name,
		}

		if !strings.Contains(model.ID, ":") {
			entry.ID += ":turbo"
		}

		if model.ContextLength > 0 {
			entry.Limit = &opencodeModelLimit{
				Context: model.ContextLength,
				Output:  model.ContextLength,
			}
		}

		if model.Opencode.Family != "" {
			entry.Family = model.Opencode.Family
		}

		if len(model.Opencode.Variants) > 0 {
			entry.Variants = model.Opencode.Variants
		}

		models[model.Name] = entry
	}

	output, err := json.MarshalIndent(models, "", "\t")
	if err != nil {
		fmt.Println("failed to encode output:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
