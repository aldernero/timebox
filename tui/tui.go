package tui

import (
	"fmt"
	"github.com/aldernero/timebox/tui/addui"
	"github.com/aldernero/timebox/tui/constants"
	"github.com/aldernero/timebox/tui/deleteui"
	"github.com/aldernero/timebox/tui/tableui"
	"github.com/aldernero/timebox/util"
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
	boxSummary sessionState = iota
	boxView
	boxAdd
	boxEdit
	boxDelete
	timeline
	timeAdd
	timeEdit
	timeDelete
)

type action int

const (
	none action = iota
	addItem
	editItem
	deleteItem
	reload
	awaitPrompt
)

type MainModel struct {
	state     sessionState
	prevState sessionState
	action    action
	period    util.TimePeriod
	status    string
	timebox   *util.TimeBox
	tbl       tableui.Model
	addPrompt addui.Model
	//editPrompt   editui.Model
	deletePrompt deleteui.Model
	windowSize   tea.WindowSizeMsg
}

func New(tb *util.TimeBox) MainModel {
	model := MainModel{state: boxSummary, timebox: tb, tbl: tableui.NewBoxSummary(tb, util.Week)}
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

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.prevState = m.state
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.action = none
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyLeft:
			m.period.Previous()
			m.action = reload
		case tea.KeyRight:
			m.period.Next()
			m.action = reload
		case tea.KeyEnter:
			if m.state == boxSummary {
				m.state = boxView
				m.action = reload
			}
		case tea.KeyEsc:
			if m.state == boxView {
				m.state = boxSummary
				m.action = reload
			}
		case tea.KeyCtrlD:
			switch m.state {
			case boxSummary:
				m.state = boxDelete
				m.action = deleteItem
			case boxView:
				m.state = timeDelete
				m.action = deleteItem
			case timeline:
				m.state = timeDelete
				m.action = deleteItem
			}
		case tea.KeyRunes:
			switch string(msg.Runes) {
			case "b":
				if m.state == timeline {
					m.state = boxSummary
					if m.state != m.prevState {
						m.action = reload
					} else {
						m.action = none
					}
				}
			case "t":
				if m.state == boxSummary {
					m.state = timeline
					if m.state != m.prevState {
						m.action = reload
					} else {
						m.action = none
					}
				}
			case "a":
				switch m.state {
				case boxSummary:
					m.state = boxAdd
					m.action = addItem
				case boxView, timeline:
					m.state = timeAdd
					m.action = addItem
				}
			case "e":
				switch m.state {
				case boxSummary:
					m.state = boxEdit
					m.action = editItem
				case boxView, timeline:
					m.state = timeEdit
					m.action = editItem
				}
			}
		}
	}
	switch m.state {
	case boxSummary:
		if m.prevState == boxSummary && m.action == none {
			var newModel tea.Model
			newModel, cmd = m.tbl.Update(msg)
			m.tbl = newModel.(tableui.Model)
		}
		if m.action == reload {
			m.tbl = tableui.NewBoxSummary(m.timebox, m.period.Period)
		}
	case boxView:
		if m.action == reload {
			boxName := m.tbl.GetSelectedBoxName()
			m.tbl = tableui.NewBoxView(m.timebox, m.period.Period, boxName)
		}
	case boxAdd:
		switch m.action {
		case addItem:
			m.addPrompt = addui.AddBox()
			m.action = awaitPrompt
		case awaitPrompt:
			switch m.addPrompt.State {
			case util.WasCancelled:
				m.state = boxSummary
				m.action = reload
			case util.HasResult:
				result := m.addPrompt.Result
				err := m.timebox.AddBox(result.Box())
				if err != nil {
					m.status = err.Error()
				}
				m.state = boxSummary
				m.action = reload
			case util.InUse:
				var newModel tea.Model
				newModel, cmd = m.addPrompt.Update(msg)
				m.addPrompt = newModel.(addui.Model)
			}
		}
	case boxEdit:
		break
	case boxDelete:
		break
	case timeline:
		if m.action == reload {
			m.tbl = tableui.NewTimeline(m.timebox, m.period.Period)
		}
	case timeAdd:
		if m.action == addItem {
			m.addPrompt = addui.AddSpan(m.tbl.GetSelectedBoxName())
		}
	case timeEdit:
		break
	case timeDelete:
		break
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var help string
	help = ""
	top := lipgloss.JoinHorizontal(lipgloss.Center, topView(), help)
	session := sessionView()
	lw := lipgloss.Width(session)
	rw := constants.TUIWidth - lw - 1
	period := lipgloss.NewStyle().Width(rw).Align(lipgloss.Right).Render(m.periodView())
	bottom := lipgloss.JoinHorizontal(lipgloss.Top, session, period)
	var middle string
	switch m.state {
	case boxSummary, boxView, timeline:
		middle = m.tbl.View()
	case boxAdd:
		middle = m.addPrompt.View()
	}
	return lipgloss.JoinVertical(lipgloss.Left, top, middle, bottom)
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
			lipgloss.Top, constants.SessionShortcutStyle.Render("<t> "), constants.SessionTextStyle.Render("Timeline")),
	)
	return lipgloss.NewStyle().PaddingLeft(1).Render(view)
}
