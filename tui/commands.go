package tui

import tea "github.com/charmbracelet/bubbletea"

type cancelMsg struct {
}

type reloadWithStatusMsg struct {
	status string
}

func reloadWithStatusCmd(status string) tea.Cmd {
	return func() tea.Msg {
		return reloadWithStatusMsg{status}
	}
}

func cancelCmd() tea.Cmd {
	return func() tea.Msg {
		return cancelMsg{}
	}
}
