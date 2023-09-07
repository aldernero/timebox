package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type DeletePrompt struct {
	txt       string
	HasAnswer bool
	Answer    bool
}

func NewDeletePrompt(txt string) DeletePrompt {
	return DeletePrompt{txt: txt}
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
	return DeleteStyle(
		fmt.Sprintf(
			"Are you sure you want to delete this (y/n)\n\n%s",
			m.txt,
		),
	)
}
