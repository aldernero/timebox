package tableui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg struct {
	error
}

func (m summaryui.Model) handleNav() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}

func (m summaryui.Model) handleAdd() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}

func (m summaryui.Model) handleEdit() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}

func (m summaryui.Model) handleDelete() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}
