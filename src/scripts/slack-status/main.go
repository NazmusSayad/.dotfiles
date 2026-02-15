package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	slack "dotfiles/src/helpers/slack"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/logrusorgru/aurora/v4"
)

func main() {
	initialStatus := readSlackStatus()
	renderSlackStatus("Current Slack Status", initialStatus)
	fmt.Println()

	items := []list.Item{
		statusItem{"Always", "Start Slack on login"},
		statusItem{"Work Time", "Start only during office hours"},
		statusItem{"Disabled", "Do not start Slack automatically"},
	}

	delegate := inlineDelegate{
		titleStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		selectedStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		descStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}

	tuiList := list.New(items, delegate, 40, 0)
	tuiList.Styles.TitleBar = lipgloss.NewStyle()

	switch initialStatus {
	case slack.SlackStatusAlways:
		tuiList.Select(0)
	case slack.SlackStatusWorkTime:
		tuiList.Select(1)
	case slack.SlackStatusDisabled:
		tuiList.Select(2)
	}

	tuiList.SetShowHelp(false)
	tuiList.SetShowTitle(false)
	tuiList.SetShowStatusBar(false)
	tuiList.SetShowPagination(false)
	tuiList.SetFilteringEnabled(false)

	m := model{list: tuiList}
	p := tea.NewProgram(m)

	final, err := p.Run()
	if err != nil {
		return
	}

	if fm, ok := final.(model); ok && fm.choice != nil {
		fmt.Println()
		writeSlackStatus(*fm.choice)
	}
}

func getSlackStatusFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".slack-status")
}

func readSlackStatus() slack.SlackStatus {
	data, err := os.ReadFile(getSlackStatusFilePath())
	if err != nil {
		return slack.SlackStatusWorkTime
	}

	status := strings.TrimSpace(string(data))
	return slack.SlackStatus(status)
}

func writeSlackStatus(status slack.SlackStatus) {
	renderSlackStatus("Updating slack status to", status)
	os.WriteFile(getSlackStatusFilePath(), []byte(status), 0644)
	slack.SlackLaunch(status)
}

func renderSlackStatus(label string, status slack.SlackStatus) {
	switch status {
	case slack.SlackStatusAlways:
		fmt.Println("> " + label + ": " + aurora.Green("Always On").String())
	case slack.SlackStatusWorkTime:
		fmt.Println("> " + label + ": " + aurora.Yellow("Work Time").String())
	case slack.SlackStatusDisabled:
		fmt.Println("> " + label + ": " + aurora.Red("Disabled").String())
	}
}

type inlineDelegate struct {
	titleStyle    lipgloss.Style
	selectedStyle lipgloss.Style
	descStyle     lipgloss.Style
}

func (d inlineDelegate) Height() int                               { return 1 }
func (d inlineDelegate) Spacing() int                              { return 0 }
func (d inlineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d inlineDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(statusItem)
	if !ok || m.Width() <= 0 {
		return
	}
	title := i.Title()
	descSuffix := ""
	if i.description != "" {
		descSuffix = "  " + i.description
	}
	textwidth := m.Width() - 2
	full := title + descSuffix
	full = ansi.Truncate(full, textwidth, "")
	prefix := "  "
	if index == m.Index() {
		prefix = d.selectedStyle.Render("✓ ")
	}
	style := d.titleStyle
	if index == m.Index() {
		style = d.selectedStyle
	}
	var line string
	if len(descSuffix) > 0 && len(full) > len(title) {
		line = prefix + style.Render(title) + d.descStyle.Render(full[len(title):])
	} else {
		line = prefix + style.Render(title)
	}
	fmt.Fprint(w, line)
}

type statusItem struct {
	title       string
	description string
}

func (i statusItem) Title() string       { return i.title }
func (i statusItem) Description() string { return i.description }
func (i statusItem) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	choice *slack.SlackStatus
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if item, ok := m.list.SelectedItem().(statusItem); ok {
				var s slack.SlackStatus
				switch item.title {
				case "Always":
					s = slack.SlackStatusAlways
				case "Work Time":
					s = slack.SlackStatusWorkTime
				case "Disabled":
					s = slack.SlackStatusDisabled
				}
				m.choice = &s
				return m, tea.Quit
			}
		}
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, min(msg.Height, 4))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
