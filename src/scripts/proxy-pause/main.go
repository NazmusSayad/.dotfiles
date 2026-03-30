package main

import (
	"os"
	"strconv"
	"time"

	"dotfiles/src/helpers"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/logrusorgru/aurora/v4"
)

const TOTAL_WAIT_SECONDS = 5

type tickMsg time.Time

type model struct {
	remaining int
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tickMsg:
		m.remaining--
		if m.remaining <= 0 {
			return m, tea.Quit
		}

		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	return aurora.Faint("Press any key to exit, or wait " + strconv.Itoa(m.remaining) + " seconds...").String()
}

func pressAnyKeyOrWaitToExit() {
	p := tea.NewProgram(model{remaining: TOTAL_WAIT_SECONDS})
	if _, err := p.Run(); err != nil {
		panic("failed to run countdown")
	}
}

func main() {
	helpers.ExecNativeCommand(os.Args[1:])
	pressAnyKeyOrWaitToExit()
}
