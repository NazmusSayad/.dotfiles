package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	helpers "dotfiles/src/helpers"
	"dotfiles/src/utils"

	"github.com/logrusorgru/aurora/v4"
	_ "modernc.org/sqlite"
)

var TARGETS = []string{
	`$APPDATA/Code/User/globalStorage/state.vscdb`,
	`$APPDATA/Cursor/User/globalStorage/state.vscdb`,
}

type StateConfig map[string]json.RawMessage

func main() {
	config := helpers.ReadConfig[StateConfig]("@/config/vscode/state.jsonc")
	keys := make([]string, 0, len(config))
	for k := range config {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return
	}

	var source string
	for _, file := range TARGETS {
		fullPath := helpers.ResolvePath(file)
		if utils.IsFileExists(fullPath) {
			source = fullPath
			break
		}
	}
	if source == "" {
		fmt.Fprintf(os.Stderr, "UNEXPECTED: no state.vscdb found in targets\n")
		os.Exit(1)
	}

	pulled := pullState(source, keys)
	outPath := helpers.ResolvePath("@/config/vscode/state.jsonc")
	if err := writeState(outPath, pulled); err != nil {
		fmt.Fprintf(os.Stderr, "UNEXPECTED: write: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(aurora.Green("Saved"), aurora.Faint(outPath))
}

func pullState(file string, keys []string) StateConfig {
	db, err := sql.Open("sqlite", "file:"+file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "UNEXPECTED: open db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "UNEXPECTED: db: %v\n", err)
		os.Exit(1)
	}

	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='ItemTable'").Scan(&name)
	if err == sql.ErrNoRows {
		fmt.Fprintf(os.Stderr, "UNEXPECTED: ItemTable not found\n")
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "UNEXPECTED: db: %v\n", err)
		os.Exit(1)
	}

	out := make(StateConfig, len(keys))
	for _, key := range keys {
		var value string
		err := db.QueryRow("SELECT value FROM ItemTable WHERE key = ?", key).Scan(&value)
		if err == sql.ErrNoRows {
			fmt.Fprintf(os.Stderr, "WARN: key not in DB: %s\n", key)
			continue
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "UNEXPECTED: %s: %v\n", key, err)
			os.Exit(1)
		}
		fmt.Println(aurora.Cyan("> Pulling"), aurora.Faint(key))
		out[key] = valueToRaw(value)
	}
	return out
}

func valueToRaw(s string) json.RawMessage {
	var parsed interface{}
	if err := json.Unmarshal([]byte(s), &parsed); err == nil {
		raw, _ := json.Marshal(parsed)
		return raw
	}
	raw, _ := json.Marshal(s)
	return raw
}

func writeState(path string, config StateConfig) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}
