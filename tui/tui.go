package tui

import (
	"fmt"
	"github.com/aldernero/timebox/tui/summaryui"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type sessionState int

const (
	summaryView sessionState = iota
	detailView
)

type MainModel struct {
	state      sessionState
	summary    summaryui.Model
	windowSize tea.WindowSizeMsg
}

func New() MainModel {
	model := MainModel{state: summaryView}
	return model
}

func StartTea() {
	m := New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MainModel) View() string {
	return m.summary.View()
}
