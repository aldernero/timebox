package tableui

import (
	"fmt"
	"github.com/aldernero/timebox/tui/constants"
	"github.com/aldernero/timebox/util"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"os"
	"strings"
	"time"
)

type mode int

const (
	nav mode = iota
	add
	edit
	delete
)

type keymap struct {
	Add    key.Binding
	Remove key.Binding
	Next   key.Binding
	Prev   key.Binding
	Enter  key.Binding
	Back   key.Binding
	Quit   key.Binding
}

// Keymap reusable key mappings shared across models
var Keymap = keymap{
	Add: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "add"),
	),
	Remove: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "remove"),
	),
	Next: key.NewBinding(
		key.WithKeys("tab"),
	),
	Prev: key.NewBinding(
		key.WithKeys("shift+tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctlr+c", "q"),
		key.WithHelp("q", "back"),
	),
}

type inputFields int

const (
	inputNameField inputFields = iota
	inputMinTimeField
	inputMaxTimeField
	inputCancelButton
	inputSubmitButton
)

type Model struct {
	mode        mode
	focus       int
	when        util.Period
	inputs      []textinput.Model
	inputStatus string
	sumTable    table.Model
	timebox     *util.TimeBox
}

func New(tb *util.TimeBox, p util.Period) Model {
	var m Model
	m.when = p
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
	m.timebox = tb
	m.sumTable = makeTable(m.timebox, p)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.mode {
	case nav:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case msg.String() == "ctrl+c":
				return m, tea.Quit
			case msg.String() == "a":
				m.mode = add
			case msg.String() == "e":
				m.mode = edit
			case msg.String() == "ctrl+d":
				m.mode = delete
			case key.Matches(msg, Keymap.Add):
				m.mode = add
			}
		}
		newTable, newCmd := m.sumTable.Update(msg)
		m.sumTable = newTable
		cmd = newCmd
	case add:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, Keymap.Back):
				m.resetInputs()
				m.mode = nav
			case key.Matches(msg, Keymap.Next):
				m.focus++
				if m.focus > int(inputSubmitButton) {
					m.focus = int(inputNameField)
				}
			case key.Matches(msg, Keymap.Prev):
				m.focus--
				if m.focus < int(inputNameField) {
					m.focus = int(inputSubmitButton)
				}
			case key.Matches(msg, Keymap.Enter):
				switch inputFields(m.focus) {
				case inputNameField, inputMinTimeField, inputMaxTimeField:
					m.focus++
				case inputCancelButton:
					m.resetInputs()
					m.mode = nav
				case inputSubmitButton:
					box, err := m.validateInputs()
					if err != nil {
						m.inputs[inputNameField].Reset()
						m.inputs[inputMinTimeField].Reset()
						m.inputs[inputMaxTimeField].Reset()
						m.focus = 0
						m.inputStatus = fmt.Sprintf("invalid inputs: %v", err)
						break
					}
					err = m.timebox.AddBox(box)
					if err != nil {
						m.inputs[inputNameField].Reset()
						m.inputs[inputMinTimeField].Reset()
						m.inputs[inputMaxTimeField].Reset()
						m.focus = 0
						m.inputStatus = fmt.Sprintf("error adding to db: %v", err)
						break
					}
					//newTable, newCmd := m.sumTable.Update(msg)
					m.sumTable = makeTable(m.timebox, m.when)
					//cmd = newCmd
					m.resetInputs()
					m.mode = nav
				}
			}
		}
		cmds = append(cmds, m.updateInputs()...)
		for i := 0; i < len(m.inputs); i++ {
			newModel, cmd := m.inputs[i].Update(msg)
			m.inputs[i] = newModel
			cmds = append(cmds, cmd)
		}
	case edit:
		cmd := m.handleEdit()
		cmds = append(cmds, cmd)
	case delete:
		cmd := m.handleDelete()
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.mode {
	case add:
		return m.inputView()
	default:
		return m.sumTable.View()
	}
}

func (m Model) Help() string {
	var view string
	switch m.mode {
	case nav:
		view = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.JoinHorizontal(
				lipgloss.Top, constants.ShortcutStyle.Render("<a>     "), constants.HelpStyle.Render("Add")),
			lipgloss.JoinHorizontal(
				lipgloss.Top, constants.ShortcutStyle.Render("<ctrl-d>"), constants.HelpStyle.Render("Delete")),
			lipgloss.JoinHorizontal(
				lipgloss.Top, constants.ShortcutStyle.Render("<e>     "), constants.HelpStyle.Render("Edit")),
			lipgloss.JoinHorizontal(
				lipgloss.Top, constants.ShortcutStyle.Render("<ctrl-c>"), constants.HelpStyle.Render("Quit")),
		)
	default:
		view = lipgloss.JoinHorizontal(
			lipgloss.Top, constants.ShortcutStyle.Render("<ctrl-c>"), constants.HelpStyle.Render("Quit"))
	}
	return constants.HelpBlockStyle.Render(view)
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
	if m.focus == len(m.inputs) {
		cancelButton = &constants.FocusedStyle
	}
	submitButton := &constants.BlurredStyle
	if m.focus == len(m.inputs)+1 {
		submitButton = &constants.FocusedStyle
	}
	_, err := fmt.Fprintf(
		&b,
		"\n\n%s  %s\n\n%s",
		cancelButton.Render("[ Cancel ]"),
		submitButton.Render("[ Submit ]"),
		constants.ErrStyle(m.inputStatus),
	)
	if err != nil {
		fmt.Printf("Error formatting input string: %v\n", err)
		os.Exit(1)
	}

	return constants.PromptStyle.Render(constants.InputStyle(b.String()))
}

func (m *Model) updateInputs() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focus {
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
	m.inputs[inputNameField].Reset()
	m.inputs[inputMinTimeField].Reset()
	m.inputs[inputMaxTimeField].Reset()
	m.focus = 0
	m.inputStatus = ""
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
	box = util.Box{name, minTime, maxTime}
	return box, nil
}
