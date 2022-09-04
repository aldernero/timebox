package summaryui

import tea "github.com/charmbracelet/bubbletea"

type errMsg struct {
	error
}

func (m Model) handleNav() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}

func (m Model) handleAdd() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}

func (m Model) handleEdit() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}

func (m Model) handleDelete() tea.Cmd {
	return func() tea.Msg {
		return errMsg{nil}
	}
}