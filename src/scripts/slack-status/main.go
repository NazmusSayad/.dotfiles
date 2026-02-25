package main

import (
	"fmt"
	"io"

	slack "dotfiles/src/helpers/slack"
	"dotfiles/src/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/logrusorgru/aurora/v4"
)

func formatOfficeTime() string {
	config := slack.ReadSlackConfig()
	timeString := utils.HourToAmPm(config.OfficeTimeStart) + "-" + utils.HourToAmPm(config.OfficeTimeFinish)

	return timeString
}

func main() {
	initialStatus := slack.GetSlackStartupConfig()
	renderSlackStatus("Current Slack Status", initialStatus)
	fmt.Println()

	items := []list.Item{
		statusItem{"Always Active", "Always keep slack running", slack.SlackStatusAlways},
		statusItem{"Office Hours", "Office Hours: " + formatOfficeTime(), slack.SlackStatusWorkTime},
		statusItem{"Never", "Keep slack disabled", slack.SlackStatusDisabled},
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
		panic("failed to run program")
	}

	fm, ok := final.(model)
	if !ok {
		panic("failed to cast final to model")
	}

	if fm.choice == nil || *fm.choice == initialStatus {
		return
	}

	renderSlackStatus("Updating slack status to", *fm.choice)
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

	slack.WriteSlackStartupConfig(status)
	slack.SlackLaunch(status)
}

type inlineDelegate struct {
	titleStyle    lipgloss.Style
	selectedStyle lipgloss.Style
	descStyle     lipgloss.Style
}

type statusItem struct {
	title       string
	description string
	status      slack.SlackStatus
}

type model struct {
	list   list.Model
	choice *slack.SlackStatus
}

func (d inlineDelegate) Height() int                               { return 1 }
func (d inlineDelegate) Spacing() int                              { return 0 }
func (d inlineDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d inlineDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(statusItem)
	if !ok || m.Width() <= 0 {
		return
	}
	selected := index == m.Index()
	prefix := "  "
	style := d.titleStyle
	if selected {
		prefix = d.selectedStyle.Render("✓ ")
		style = d.selectedStyle
	}
	descSuffix := ""
	if i.description != "" {
		descSuffix = "  " + i.description
	}
	full := ansi.Truncate(i.title+descSuffix, m.Width()-2, "")
	if descSuffix != "" && len(full) > len(i.title) {
		fmt.Fprint(w, prefix+style.Render(i.title)+d.descStyle.Render(full[len(i.title):]))
	} else {
		fmt.Fprint(w, prefix+style.Render(i.title))
	}
}

func (i statusItem) Title() string       { return i.title }
func (i statusItem) Description() string { return i.description }
func (i statusItem) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if k == "enter" {
			if item, ok := m.list.SelectedItem().(statusItem); ok {
				m.choice = &item.status
				return m, tea.Quit
			}
		}
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, min(msg.Height, 5))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
