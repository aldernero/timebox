package tui

import (
	"fmt"
	util2 "github.com/aldernero/timebox/pkg/util"
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

type AddPrompt struct {
	mode         promptType
	State        util2.PromptState
	focusedField inputFields
	editMode     bool
	inputs       []textinput.Model
	status       string
	Result       util2.InputResult
}

func AddBox() AddPrompt {
	var m AddPrompt
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.PromptStyle = NoStyle
		t.Cursor.Style = NoStyle
		t.CharLimit = 30
		switch i {
		case 0:
			t.Prompt = "Name  > "
			t.Placeholder = "Box Name"
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
		case 1:
			t.Prompt = "Min   > "
			t.Placeholder = "Weekly Min (e.g. 1h30m)"
			t.CharLimit = 30
		case 2:
			t.Prompt = "Max   > "
			t.Placeholder = "Weekly Max (e.g. 4h)"
			t.CharLimit = 30
		}
		m.inputs[i] = t
	}
	m.State = util2.InUse
	return m
}

func EditBox(box util2.Box) AddPrompt {
	var m AddPrompt
	m.editMode = true
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.PromptStyle = NoStyle
		t.Cursor.Style = NoStyle
		t.CharLimit = 30
		switch i {
		case 0:
			t.Prompt = "Name > "
			t.SetValue(box.Name)
		case 1:
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
			t.Prompt = "Min  > "
			t.SetValue(box.MinTime.String())
			t.CharLimit = 30
		case 2:
			t.Prompt = "Max  > "
			t.SetValue(box.MaxTime.String())
			t.CharLimit = 30
		}
		m.inputs[i] = t
	}
	m.focusedField = minField
	m.State = util2.InUse
	return m
}

func AddSpan(boxName string) AddPrompt {
	var m AddPrompt
	m.mode = spanInput
	m.inputs = make([]textinput.Model, 3)
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.PromptStyle = NoStyle
		t.CursorStyle = NoStyle
		t.CharLimit = 30
		switch i {
		case 0:
			t.Prompt = "Box  > "
			t.SetValue(boxName)
		case 1:
			t.Prompt = "Start > "
			t.Placeholder = "Start (e.g. 2023-04-01 17:35:00)"
			t.CharLimit = 30
			t.Focus()
			t.PromptStyle = FocusedStyle
			t.TextStyle = FocusedStyle
		case 2:
			t.Prompt = "End   > "
			t.Placeholder = "End (e.g. 2023-04-01 19:02:30)"
			t.CharLimit = 30
		}
		m.inputs[i] = t
	}
	m.focusedField = minField
	m.State = util2.InUse
	return m
}

func (m AddPrompt) Init() tea.Cmd {
	return nil
}

func (m AddPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var checkInput bool
	var baseField int
	if m.editMode {
		baseField = 1
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focusedField++
			if m.focusedField > submitButton {
				m.focusedField = inputFields(baseField)
			}
		case "shift+tab":
			m.focusedField--
			if int(m.focusedField) < baseField {
				m.focusedField = submitButton
			}
		case "enter":
			switch m.focusedField {
			case cancelButton:
				m.State = util2.WasCancelled
				return m, nil
			case submitButton:
				checkInput = true
			}
		}
	}
	if checkInput {
		switch m.mode {
		case boxInput:
			box, err := m.validateBoxInputs()
			if err != nil {
				m.status = err.Error()
			}
			m.Result = util2.NewInputResultBox(box)
			m.State = util2.HasResult
			return m, nil
		case spanInput:
			span, err := m.validateSpanInputs()
			if err != nil {
				m.status = err.Error()
			}
			m.Result = util2.NewInputResultSpan(span)
			m.State = util2.HasResult
			return m, nil
		}
	}
	cmds = append(cmds, m.updateInputs()...)
	for i := 0; i < len(m.inputs); i++ {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m AddPrompt) View() string {
	return m.inputView()
}

func (m AddPrompt) inputView() string {
	var b strings.Builder
	var title string
	switch m.mode {
	case boxInput:
		if m.editMode {
			title = "Edit Box"
		} else {
			title = "New Box"
		}
	case spanInput:
		title = "New Timespan"
	}
	b.WriteString(InputTitleStyle.Render(title) + "\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	cancelButton := &BlurredStyle
	if int(m.focusedField) == len(m.inputs) {
		cancelButton = &FocusedStyle
	}
	submitButton := &BlurredStyle
	if int(m.focusedField) == len(m.inputs)+1 {
		submitButton = &FocusedStyle
	}
	_, err := fmt.Fprintf(
		&b,
		"\n\n%s  %s\n\n%s",
		cancelButton.Render("[ Cancel ]"),
		submitButton.Render("[ Submit ]"),
		ErrStyle(m.status),
	)
	if err != nil {
		fmt.Printf("Error formatting input string: %v\n", err)
		os.Exit(1)
	}

	return PromptStyle.Render(InputStyle(b.String()))
}

func (m AddPrompt) updateInputs() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == int(m.focusedField) {
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

func (m AddPrompt) resetInputs() {
	m.inputs[nameField].Reset()
	m.inputs[minField].Reset()
	m.inputs[maxField].Reset()
	m.focusedField = nameField
	m.status = ""
}

func (m AddPrompt) validateBoxInputs() (util2.Box, error) {
	var box util2.Box
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
	box = util2.Box{Name: name, MinTime: minTime, MaxTime: maxTime}
	return box, nil
}

func (m AddPrompt) validateSpanInputs() (util2.Span, error) {
	var span util2.Span
	name := m.inputs[0].Value()
	min := m.inputs[1].Value()
	max := m.inputs[2].Value()
	if name == "" || min == "" || max == "" {
		return span, fmt.Errorf("empty fields")
	}
	minTime, err := util2.ParseTime(min)
	if err != nil {
		return span, fmt.Errorf("invalid duration: %v", err)
	}
	maxTime, err := util2.ParseTime(max)
	if err != nil {
		return span, fmt.Errorf("invalid duration: %v", err)
	}
	span = util2.Span{
		Start: minTime,
		End:   maxTime,
		Box:   name,
	}
	return span, nil
}
