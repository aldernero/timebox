package tui

import tea "github.com/charmbracelet/bubbletea"

type DeletePrompt struct {
	text      string
	HasAnswer bool
	Answer    bool
}

func NewDeletePrompt(txt string) DeletePrompt {
	return DeletePrompt{text: txt}
}

func (m DeletePrompt) Init() tea.Cmd {
	return nil
}

func (m DeletePrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			m.HasAnswer = true
			m.Answer = false
		case "y":
			m.HasAnswer = true
			m.Answer = true
		}
	}
	return m, cmd
}

func (m DeletePrompt) View() string {
	return "Are you sure you want to delete " + m.text + "? (y/n)"
}
