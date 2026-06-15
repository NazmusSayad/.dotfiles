package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	localDir := resolveDotfilesLocalDir()
	targetPath := resolveClaudeCredentialsPath()

	accounts := readAccounts(localDir, targetPath)
	if len(accounts) == 0 {
		fmt.Println(aurora.Red("No *.credentials.json files found in .local"))
		return
	}

	currentAccount := "unknown"
	for _, a := range accounts {
		if a.isCurrent {
			currentAccount = a.name
		}
	}
	fmt.Println("> Current Claude account: " + aurora.Green(currentAccount).String())

	lastAccount := "unknown"
	if data, err := os.ReadFile(filepath.Join(localDir, ".claude-last-account")); err == nil {
		if s := strings.TrimSpace(string(data)); s != "" {
			lastAccount = s
		}
	}
	fmt.Println("> Last used account: " + aurora.Green(lastAccount).String())
	fmt.Println()

	choice := selectAccount(accounts)
	if choice == nil {
		return
	}

	data, err := os.ReadFile(targetPath)
	if err != nil {
		fmt.Println(aurora.Red("Failed to read claude credentials: " + targetPath))
		return
	}
	err = os.WriteFile(choice.src, data, 0o600)
	if err != nil {
		fmt.Println(aurora.Red("Failed to write to local: " + choice.src))
		return
	}
	fmt.Println("> Pulled Claude credentials to " + aurora.Green(choice.name).String())

	_ = os.WriteFile(filepath.Join(localDir, ".claude-last-account"), []byte(choice.name), 0o600)
}

func resolveDotfilesLocalDir() string {
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

func resolveClaudeCredentialsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".claude", ".credentials.json")
}

type claudeCreds struct {
	ClaudeAiOauth struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	} `json:"claudeAiOauth"`
}

func readTokens(path string) (access, refresh string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", ""
	}
	var c claudeCreds
	if json.Unmarshal(data, &c) != nil {
		return "", ""
	}
	return c.ClaudeAiOauth.AccessToken, c.ClaudeAiOauth.RefreshToken
}

type account struct {
	name      string
	src       string
	isCurrent bool
}

func readAccounts(dir, targetPath string) []account {
	curAccess, curRefresh := readTokens(targetPath)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	accounts := []account{}
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

		accounts = append(accounts, account{
			name:      strings.TrimSuffix(e.Name(), ".credentials.json"),
			src:       src,
			isCurrent: isCurrent,
		})
	}
	return accounts
}

func selectAccount(accounts []account) *account {
	items := make([]list.Item, len(accounts))
	current := -1
	for i, a := range accounts {
		if a.isCurrent {
			current = i
		}
		items[i] = accountItem{title: a.name, account: a}
	}

	delegate := inlineDelegate{
		titleStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		selectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	}

	tuiList := list.New(items, delegate, 40, 0)
	tuiList.Styles.TitleBar = lipgloss.NewStyle()
	tuiList.SetShowHelp(false)
	tuiList.SetShowTitle(false)
	tuiList.SetShowStatusBar(false)
	tuiList.SetShowPagination(false)
	tuiList.SetFilteringEnabled(false)
	if current >= 0 {
		tuiList.Select(current)
	}

	final, err := tea.NewProgram(model{list: tuiList}).Run()
	if err != nil {
		panic("failed to run program")
	}

	fm, ok := final.(model)
	if !ok {
		panic("failed to cast final to model")
	}
	return fm.choice
}

type inlineDelegate struct {
	titleStyle    lipgloss.Style
	selectedStyle lipgloss.Style
}

type accountItem struct {
	title   string
	account account
}

type model struct {
	list   list.Model
	choice *account
}

func (d inlineDelegate) Height() int                               { return 1 }
func (d inlineDelegate) Spacing() int                              { return 0 }
func (d inlineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d inlineDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(accountItem)
	if !ok || m.Width() <= 0 {
		return
	}

	prefix := "  "
	style := d.titleStyle
	if index == m.Index() {
		prefix = d.selectedStyle.Render("✓ ")
		style = d.selectedStyle
	}

	full := ansi.Truncate(i.title, m.Width()-2, "")
	fmt.Fprint(w, prefix+style.Render(full))
}

func (i accountItem) Title() string       { return i.title }
func (i accountItem) Description() string { return "" }
func (i accountItem) FilterValue() string { return i.account.name }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "enter" {
			if item, ok := m.list.SelectedItem().(accountItem); ok {
				m.choice = &item.account
				return m, tea.Quit
			}
		}
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, min(msg.Height, len(m.list.Items()))+2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
