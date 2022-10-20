package tui

import (
	"fmt"
	"github.com/aldernero/timebox/tui/constants"
	"github.com/aldernero/timebox/tui/summaryui"
	"github.com/aldernero/timebox/util"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

const (
	logoFile = "assets/timebox.txt"
)

type sessionState int

const (
	summaryView sessionState = iota
	detailView
)

type MainModel struct {
	state      sessionState
	period     util.TimePeriod
	timebox    *util.TimeBox
	summary    summaryui.Model
	windowSize tea.WindowSizeMsg
}

func New(tb *util.TimeBox) MainModel {
	model := MainModel{state: summaryView, timebox: tb, summary: summaryui.New(tb, util.Week)}
	return model
}

func StartTea(tb *util.TimeBox) {
	m := New(tb)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

type keymap struct {
	left  key.Binding
	right key.Binding
	boxes key.Binding
	spans key.Binding
	quit  key.Binding
}

var Keymap = keymap{
	left: key.NewBinding(
		key.WithKeys("left"),
	),
	right: key.NewBinding(
		key.WithKeys("right"),
	),
	boxes: key.NewBinding(
		key.WithKeys("b"),
	),
	spans: key.NewBinding(
		key.WithKeys("s"),
	),
	quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
	),
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.left):
			m.period.Previous()
			m.summary = summaryui.New(m.timebox, m.period.Period)
		case key.Matches(msg, Keymap.right):
			m.period.Next()
			m.summary = summaryui.New(m.timebox, m.period.Period)
		case key.Matches(msg, Keymap.boxes):
			m.state = summaryView
		case key.Matches(msg, Keymap.spans):
			m.state = detailView
		case key.Matches(msg, Keymap.quit):
			return m, tea.Quit
		}
	}
	switch m.state {
	case summaryView:
		newSummary, newCmd := m.summary.Update(msg)
		newModel, ok := newSummary.(summaryui.Model)
		if !ok {
			panic("couldn't perform assertion on summaryui model")
		}
		m.summary = newModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var help string
	help = m.summary.Help()
	top := lipgloss.JoinHorizontal(lipgloss.Center, topView(), help)
	session := sessionView()
	lw := lipgloss.Width(session)
	rw := constants.TUIWidth - lw - 1
	period := lipgloss.NewStyle().Width(rw).Align(lipgloss.Right).Render(m.periodView())
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, session, period)
	return lipgloss.JoinVertical(lipgloss.Left, top, m.summary.View(), bottom)
}

func loadLogo() string {
	var logo string
	buf, err := os.ReadFile(logoFile)
	if err != nil {
		fmt.Println("Error reading logo file:", err)
	}
	logo = string(buf)
	return logo
}

func topView() string {
	var view string
	logo := loadLogo()
	view = constants.LogoStyle.Render(logo)
	return view
}

func (m MainModel) periodView() string {
	var b strings.Builder
	names := m.period.Names()
	current := m.period.Current()
	for i, name := range names {
		if i == current {
			b.WriteString(constants.CurrentPeriodStyle.Render(name))
		} else {
			b.WriteString(constants.PeriodStyle.Render(name))
		}
	}
	return constants.PeriodPickerStyle.Render(b.String())
}

func sessionView() string {
	var view string
	view = lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top, constants.SessionShortcutStyle.Render("<b> "), constants.SessionTextStyle.Render("Boxes")),
		lipgloss.JoinHorizontal(
			lipgloss.Top, constants.SessionShortcutStyle.Render("<s> "), constants.SessionTextStyle.Render("Spans")),
		lipgloss.JoinHorizontal(
			lipgloss.Top, constants.SessionShortcutStyle.Render("<t> "), constants.SessionTextStyle.Render("Timeline")),
	)
	return lipgloss.NewStyle().PaddingLeft(1).Render(view)
}
