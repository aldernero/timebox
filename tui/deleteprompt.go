package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

type DeletePrompt struct {
	category  string
	item      string
	HasAnswer bool
	Answer    bool
}

func NewDeletePrompt(category, item string) DeletePrompt {
	return DeletePrompt{category: category, item: item}
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
			"Are you sure you want to delete %s: %s? (y/n)",
			m.category,
			m.item),
	)
}
