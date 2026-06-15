package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const LastAccountFile = ".claude-last-account"

type Account struct {
	Name      string
	Path      string
	IsCurrent bool
}

type creds struct {
	ClaudeAiOauth struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"claudeAiOauth"`
}

func ResolveLocalDir() string {
	dotfiles := os.Getenv("DOTFILES_DIR")
	if dotfiles == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		cand := filepath.Join(home, ".dotfiles")
		if fi, err := os.Stat(cand); err == nil && fi.IsDir() {
			dotfiles = cand
		} else {
			panic("DOTFILES_DIR environment variable is not set")
		}
	}
	return filepath.Join(dotfiles, ".local")
}

func ResolveCredentialsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".claude", ".credentials.json")
}

func readTokens(path string) (access, refresh string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}
	var c creds
	if json.Unmarshal(data, &c) != nil {
		return "", ""
	}
	return c.ClaudeAiOauth.AccessToken, c.ClaudeAiOauth.RefreshToken
}

func ReadAccounts(dir, targetPath string) []Account {
	curAccess, curRefresh := readTokens(targetPath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	accounts := []Account{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), ".credentials.json") {
			continue
		}

		src := filepath.Join(dir, e.Name())
		access, refresh := readTokens(src)
		isCurrent := (curAccess != "" && access == curAccess) || (curRefresh != "" && refresh == curRefresh)

		accounts = append(accounts, Account{
			Name:      strings.TrimSuffix(e.Name(), ".credentials.json"),
			Path:      src,
			IsCurrent: isCurrent,
		})
	}
	return accounts
}

func CurrentAccount(accounts []Account) string {
	for _, a := range accounts {
		if a.IsCurrent {
			return a.Name
		}
	}
	return "unknown"
}

func ReadLastAccount(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, LastAccountFile))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func WriteLastAccount(dir, name string) {
	_ = os.WriteFile(filepath.Join(dir, LastAccountFile), []byte(name), 0o600)
}
