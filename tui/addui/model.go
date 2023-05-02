package addui

import (
	"fmt"
	"github.com/aldernero/timebox/tui/constants"
	"github.com/aldernero/timebox/util"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
	"time"
)

type promptType int

const (
	boxInput promptType = iota
	spanInput
)

type inputFields int

const (
	nameField inputFields = iota
	minField
	maxField
	cancelButton
	submitButton
)

type state int

const (
	normal state = iota
	focused
	disabled
)

type Model struct {
	mode           promptType
	state          state
	focusedField   inputFields
	editMode       bool
	inputs         []textinput.Model
	status         string
	nameFieldTitle string
	minFieldTitle  string
	maxFieldTitle  string
}

func AddBox() Model {
	var m Model
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.PromptStyle = constants.NoStyle
		t.CursorStyle = constants.NoStyle
		t.CharLimit = 30
		switch i {
		case 0:
			t.Placeholder = "Box Name"
			t.Focus()
			t.PromptStyle = constants.FocusedStyle
			t.TextStyle = constants.FocusedStyle
		case 1:
			t.Placeholder = "Weekly Min (e.g. 1h30m)"
			t.CharLimit = 30
		case 2:
			t.Placeholder = "Weekly Max (e.g. 4h)"
			t.CharLimit = 30
		}
		m.inputs[i] = t
	}
	return m
}

func AddSpan(boxName string) Model {
	var m Model
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.PromptStyle = constants.NoStyle
		t.CursorStyle = constants.NoStyle
		t.CharLimit = 30
		switch i {
		case 0:
			t.Placeholder = "Box Name"
			t.Focus()
			t.PromptStyle = constants.FocusedStyle
			t.TextStyle = constants.FocusedStyle
		case 1:
			t.Placeholder = "Start (e.g. 2023-04-01 17:35:00)"
			t.CharLimit = 30
		case 2:
			t.Placeholder = "End (e.g. 2023-04-01 19:02:30)"
			t.CharLimit = 30
		}
		m.inputs[i] = t
	}
	m.inputs[nameField].Prompt = boxName
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
	return ""
}

func (m Model) inputView() string {
	var b strings.Builder
	b.WriteString(constants.InputTitleStyle.Render("New Box") + "\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	cancelButton := &constants.BlurredStyle
	if int(m.focusedField) == len(m.inputs) {
		cancelButton = &constants.FocusedStyle
	}
	submitButton := &constants.BlurredStyle
	if int(m.focusedField) == len(m.inputs)+1 {
		submitButton = &constants.FocusedStyle
	}
	_, err := fmt.Fprintf(
		&b,
		"\n\n%s  %s\n\n%s",
		cancelButton.Render("[ Cancel ]"),
		submitButton.Render("[ Submit ]"),
		constants.ErrStyle(m.status),
	)
	if err != nil {
		fmt.Printf("Error formatting input string: %v\n", err)
		os.Exit(1)
	}

	return constants.PromptStyle.Render(constants.InputStyle(b.String()))
}

func (m Model) updateInputs() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == int(m.focusedField) {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = constants.FocusedStyle
			m.inputs[i].TextStyle = constants.FocusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = constants.NoStyle
		m.inputs[i].TextStyle = constants.NoStyle
	}
	return cmds
}

func (m Model) resetInputs() {
	m.inputs[nameField].Reset()
	m.inputs[minField].Reset()
	m.inputs[maxField].Reset()
	m.focusedField = nameField
	m.status = ""
}

func (m Model) validateInputs() (util.Box, error) {
	var box util.Box
	name := m.inputs[0].Value()
	min := m.inputs[1].Value()
	max := m.inputs[2].Value()
	if name == "" || min == "" || max == "" {
		return box, fmt.Errorf("empty fields")
	}
	minTime, err := time.ParseDuration(min)
	if err != nil {
		return box, fmt.Errorf("invalid duration: %v", err)
	}
	maxTime, err := time.ParseDuration(max)
	if err != nil {
		return box, fmt.Errorf("invalid duration: %v", err)
	}
	box = util.Box{Name: name, MinTime: minTime, MaxTime: maxTime}
	return box, nil
}
