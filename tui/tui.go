package tui

import (
	_ "embed"
	"fmt"
	"github.com/aldernero/timebox/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"log"
	"os"
)

//go:embed timebox.txt
var logo string

type Model struct {
	state     crudState
	view      viewMode
	currScope string
	period    util.TimePeriod
	tb        util.TimeBox
	tbl       table.Model
	addPrompt AddPrompt
	delPrompt DeletePrompt
}

func New(tb util.TimeBox) Model {
	return Model{
		state:  nav,
		view:   boxSummary,
		period: util.TimePeriod{Period: util.Week},
		tb:     tb,
		tbl:    makeBoxSummaryTable(tb, util.Week),
	}
}

func StartTea(tb util.TimeBox) {
	p := tea.NewProgram(New(tb), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// CRUD state machine
	switch m.state {
	case add: // create
		return m.updateAdd(msg)
	case nav: // read
		return m.updateNav(msg)
	case edit: // update
		return m.updateEdit(msg)
	case del: // delete
		return m.updateDel(msg)
	default:
		return m, nil
	}
}

func (m Model) View() string {
	switch m.state {
	case add: // create
		return m.addPrompt.View()
	case nav: // read
		return m.mainView()
	case edit: // update
		return m.addPrompt.View()
	case del: // delete
		return m.delPrompt.View()
	default:
		return "unknown"
	}
}

func (m Model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = nav
		}
	}
	mdl, cmd := m.addPrompt.Update(msg)
	m.addPrompt = mdl.(AddPrompt)
	switch m.addPrompt.State {
	case util.WasCancelled:
		m.state = nav
	case util.HasResult:
		switch m.view {
		case boxSummary:
			res := m.addPrompt.Result
			err := m.tb.AddBox(res.Box())
			if err != nil {
				log.Fatal(err)
			}
			m.tb = util.TimeBoxFromDB(m.tb.Fname)
			m.tbl = makeBoxSummaryTable(m.tb, util.Week)
			m.state = nav
		case boxView:
			res := m.addPrompt.Result
			err := m.tb.AddSpan(res.Span(), m.currScope)
			if err != nil {
				log.Fatal(err)
			}
			m.tb = util.TimeBoxFromDB(m.tb.Fname)
			m.tbl = makeBoxViewTable(m.tb, m.currScope, util.Week)
			m.state = nav
		}
	}
	return m, cmd
}

func (m Model) updateNav(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case reloadWithStatusMsg:
		switch m.view {
		case boxSummary:
			m.tbl = makeBoxSummaryTable(m.tb, m.period.Period)
		case boxView:
			m.tbl = makeBoxViewTable(m.tb, m.currScope, m.period.Period)
		case timeline:
			m.tbl = makeTimelineTable(m.tb, m.period.Period)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit
		case "enter":
			if m.view == boxSummary {
				m.view = boxView
				boxName := m.getSelectedBoxName()
				m.currScope = boxName
				m.tbl = makeBoxViewTable(m.tb, boxName, util.Week)
			}
		case "esc":
			if m.view == boxView || m.view == timeline {
				m.view = boxSummary
				m.tbl = makeBoxSummaryTable(m.tb, util.Week)
			}
		case "tab":
			m.period.Next()
			cmd = reloadWithStatusCmd(fmt.Sprintf("Period: %s", m.period.String()))
			return m, cmd
		case "shift+tab":
			m.period.Previous()
			cmd = reloadWithStatusCmd(fmt.Sprintf("Period: %s", m.period.String()))
			return m, cmd
		case "a":
			m.state = add
			switch m.view {
			case boxSummary:
				m.addPrompt = AddBox()
			case boxView:
				m.addPrompt = AddSpan(m.currScope)
			case timeline:
				m.addPrompt = AddSpan("")
			}
		case "d":
			m.state = del
			switch m.view {
			case boxSummary:
				boxName := m.getSelectedBoxName()
				m.delPrompt = NewDeletePrompt(boxName)
			}
		case "b":
			m.view = boxSummary
			m.tbl = makeBoxSummaryTable(m.tb, util.Week)
		case "t":
			m.view = timeline
			m.tbl = makeTimelineTable(m.tb, util.Week)
		case "e":
			m.state = edit
			box := m.getSelectedBox()
			m.addPrompt = EditBox(box)
		}
	}
	m.tbl, cmd = m.tbl.Update(msg)
	return m, cmd
}

func (m Model) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = nav
		}
	}

	mdl, cmd := m.addPrompt.Update(msg)
	m.addPrompt = mdl.(AddPrompt)
	switch m.addPrompt.State {
	case util.WasCancelled:
		m.state = nav
	case util.HasResult:
		res := m.addPrompt.Result
		err := m.tb.UpdateBox(res.Box())
		if err != nil {
			log.Fatal(err)
		}
		m.tb = util.TimeBoxFromDB(m.tb.Fname)
		m.tbl = makeBoxSummaryTable(m.tb, util.Week)
		m.state = nav
	}
	return m, cmd
}

func (m Model) updateDel(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = nav
		}
	}
	mdl, cmd := m.delPrompt.Update(msg)
	m.delPrompt = mdl.(DeletePrompt)
	if m.delPrompt.HasAnswer {
		if m.delPrompt.Answer {
			boxName := m.getSelectedBoxName()
			err := m.tb.DeleteBox(boxName)
			if err != nil {
				log.Fatal(err)
			}
			m.tb = util.TimeBoxFromDB(m.tb.Fname)
			m.tbl = makeBoxSummaryTable(m.tb, util.Week)
		}
		m.state = nav
	}
	return m, cmd
}

func (m Model) mainView() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, LogoStyle.Render(logo), m.helpString()),
		m.tbl.View(),
		m.period.View())
}

func (m Model) helpString() string {
	var result string
	switch m.view {
	case boxSummary:
		row1 := ShortcutRow([]Shortcut{addShortcut, editShortcut, deleteShortcut, quitShortcut})
		row2 := ShortcutRow([]Shortcut{enterShortcut, periodShortcut, timelineShortcut})
		result = lipgloss.NewStyle().PaddingTop(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, row1, row2))
	case boxView:
		row1 := ShortcutRow([]Shortcut{addShortcut, editShortcut, deleteShortcut, quitShortcut})
		row2 := ShortcutRow([]Shortcut{backShortcut, periodShortcut})
		result = lipgloss.NewStyle().PaddingTop(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, row1, row2))
	case timeline:
		row1 := ShortcutRow([]Shortcut{editShortcut, deleteShortcut, quitShortcut})
		row2 := ShortcutRow([]Shortcut{boxSummaryShortcut, periodShortcut, timelineShortcut})
		result = lipgloss.NewStyle().PaddingTop(1).Render(lipgloss.JoinHorizontal(lipgloss.Top, row1, row2))
	}
	return result
}

func (m Model) getSelectedBoxName() string {
	row := m.tbl.HighlightedRow()
	return row.Data[columnKeyBox].(string)
}

func (m Model) getSelectedBox() util.Box {
	boxName := m.getSelectedBoxName()
	return m.tb.Boxes[boxName]
}
