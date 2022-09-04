package summaryui

import (
	"fmt"
	"github.com/aldernero/timebox/util"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"os"
	"strings"
	"time"
)

const (
	defaultListWidth   = 28
	defaultDetailWidth = 45
	defaultInputWidth  = 22
	inputTimeFormShort = "2006-01-02"
	inputTimeFormLong  = "2006-01-02 15:04:05"
	cError             = "#CF002E"
	cItemTitleDark     = "#F5EB6D"
	cItemTitleLight    = "#F3B512"
	cItemDescDark      = "#9E9742"
	cItemDescLight     = "#FFD975"
	cTitle             = "#2389D3"
	cDetailTitle       = "#D32389"
	cPromptBorder      = "#D32389"
	cDimmedTitleDark   = "#DDDDDD"
	cDimmedTitleLight  = "#222222"
	cDimmedDescDark    = "#999999"
	cDimmedDescLight   = "#555555"
	cTextLightGray     = "#FFFDF5"
)

var AppStyle = lipgloss.NewStyle().Margin(0, 1)
var TitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(cTextLightGray)).
	Background(lipgloss.Color(cTitle)).
	Padding(0, 1)
var DetailTitleStyle = lipgloss.NewStyle().
	Width(defaultDetailWidth).
	Foreground(lipgloss.Color(cTextLightGray)).
	Background(lipgloss.Color(cDetailTitle)).
	Padding(0, 1).
	Align(lipgloss.Center)
var InputTitleStyle = lipgloss.NewStyle().
	Width(defaultInputWidth).
	Foreground(lipgloss.Color(cTextLightGray)).
	Background(lipgloss.Color(cDetailTitle)).
	Padding(0, 1).
	Align(lipgloss.Center)
var SelectedTitle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(lipgloss.AdaptiveColor{Light: cItemTitleLight, Dark: cItemTitleDark}).
	Foreground(lipgloss.AdaptiveColor{Light: cItemTitleLight, Dark: cItemTitleDark}).
	Padding(0, 0, 0, 1)
var SelectedDesc = SelectedTitle.Copy().
	Foreground(lipgloss.AdaptiveColor{Light: cItemDescLight, Dark: cItemDescDark})
var DimmedTitle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: cDimmedTitleLight, Dark: cDimmedTitleDark}).
	Padding(0, 0, 0, 2)
var DimmedDesc = DimmedTitle.Copy().
	Foreground(lipgloss.AdaptiveColor{Light: cDimmedDescDark, Dark: cDimmedDescLight})
var InputStyle = lipgloss.NewStyle().
	Margin(1, 1).
	Padding(1, 2).
	Border(lipgloss.RoundedBorder(), true, true, true, true).
	BorderForeground(lipgloss.Color(cPromptBorder)).
	Render
var DetailStyle = lipgloss.NewStyle().
	Padding(1, 2).
	Border(lipgloss.ThickBorder(), false, false, false, true).
	BorderForeground(lipgloss.AdaptiveColor{Light: cItemTitleLight, Dark: cItemTitleDark}).
	Render
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cError)).Render
var NoStyle = lipgloss.NewStyle()
var FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cPromptBorder))
var BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
var BrightTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: cDimmedTitleLight, Dark: cDimmedTitleDark}).Render
var NormalTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.AdaptiveColor{Light: cDimmedDescLight, Dark: cDimmedDescDark}).Render
var SpecialTextStyle = lipgloss.NewStyle().
	Width(defaultDetailWidth).
	Margin(0, 0, 1, 0).
	Foreground(lipgloss.AdaptiveColor{Light: cItemTitleLight, Dark: cItemTitleDark}).
	Align(lipgloss.Center).Render
var DetailsBlockLeft = lipgloss.NewStyle().
	Width(defaultDetailWidth / 2).
	Foreground(lipgloss.AdaptiveColor{Light: cDimmedTitleLight, Dark: cDimmedTitleDark}).
	Align(lipgloss.Right).
	Render
var DetailsBlockRight = lipgloss.NewStyle().
	Width(defaultDetailWidth / 2).
	Foreground(lipgloss.AdaptiveColor{Light: cDimmedDescLight, Dark: cDimmedDescDark}).
	Align(lipgloss.Left).
	Render
var HelpStyle = list.DefaultStyles().HelpStyle.Width(defaultListWidth).Height(5)

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

func New(tb *util.TimeBox) Model {
	var m Model
	m.when = util.Week
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 30
		switch i {
		case 0:
			t.Placeholder = "Box Name"
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
		case 1:
			t.Placeholder = "Min Duration"
			t.CharLimit = 30
		case 2:
			t.Placeholder = "Max Duration"
			t.CharLimit = 30
		}
		m.inputs[i] = t
	}
	m.timebox = tb
	m.sumTable = makeTable(m.timebox)
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
					m.sumTable = makeTable(m.timebox)
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

func (m Model) inputView() string {
	var b strings.Builder
	b.WriteString(InputTitleStyle.Render("New Box") + "\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	cancelButton := &BlurredStyle
	if m.focus == len(m.inputs) {
		cancelButton = &FocusedStyle
	}
	submitButton := &BlurredStyle
	if m.focus == len(m.inputs)+1 {
		submitButton = &FocusedStyle
	}
	_, err := fmt.Fprintf(
		&b,
		"\n\n%s  %s\n\n%s",
		cancelButton.Render("[ Cancel ]"),
		submitButton.Render("[ Submit ]"),
		ErrStyle(m.inputStatus),
	)
	if err != nil {
		fmt.Printf("Error formatting input string: %v\n", err)
		os.Exit(1)
	}

	return InputStyle(b.String())
}

func (m *Model) updateInputs() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focus {
			// Set focused state
			cmds[i] = m.inputs[i].Focus()
			m.inputs[i].PromptStyle = FocusedStyle
			m.inputs[i].TextStyle = FocusedStyle
			continue
		}
		// Remove focused state
		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = NoStyle
		m.inputs[i].TextStyle = NoStyle
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
