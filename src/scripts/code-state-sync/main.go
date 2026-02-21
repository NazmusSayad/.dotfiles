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
	kv := helpers.ReadConfig[StateConfig]("@/config/vscode/state.jsonc")

	for _, file := range TARGETS {
		fullPath := helpers.ResolvePath(file)

		if !utils.IsFileExists(fullPath) {
			fmt.Fprintf(os.Stderr, "UNEXPECTED: file not found: %s\n", fullPath)
			continue
		}

		syncState(fullPath, kv)
	}
}

func valueToStore(raw json.RawMessage) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return string(raw)
}

func syncState(file string, config StateConfig) {
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

	for key, raw := range config {
		fmt.Println(aurora.Green("> Syncing"), aurora.Faint(key))
		value := valueToStore(raw)
		_, err := db.Exec("INSERT OR REPLACE INTO ItemTable (key, value) VALUES (?, ?)", key, value)
		if err != nil {
			fmt.Fprintf(os.Stderr, "UNEXPECTED: %s: %v\n", key, err)
			os.Exit(1)
		}
	}
}
