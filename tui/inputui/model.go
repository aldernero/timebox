package inputui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type promptType int

const (
	boxInput promptType = iota
	spanInput
)

type InputPrompt struct {
	mode        promptType
	inputs      []textinput.Model
	inputStatus string
}

func (m InputPrompt) Init() tea.Cmd {
	return nil
}

func (m InputPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	switch m.mode {
	case boxInput:
		break
	case spanInput:
		break
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m InputPrompt) View() string {
	return ""
}
