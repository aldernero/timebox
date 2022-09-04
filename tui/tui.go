package tui

import (
	"fmt"
	"github.com/aldernero/timebox/tui/summaryui"
	"github.com/aldernero/timebox/util"
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
	timebox    *util.TimeBox
	summary    summaryui.Model
	windowSize tea.WindowSizeMsg
}

func New(tb *util.TimeBox) MainModel {
	model := MainModel{state: summaryView, timebox: tb, summary: summaryui.New(tb)}
	return model
}

func StartTea(tb *util.TimeBox) {
	m := New(tb)
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
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	}
	switch m.state {
	case summaryView:
		newSummary, newCmd := m.summary.Update(msg)
		newModel, ok := newSummary.(summaryui.Model)
		if !ok {
			panic("couldn't perform assertion on summaryui model")
		}
		m.summary = newModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	return m.summary.View()
}
