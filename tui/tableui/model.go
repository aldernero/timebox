package tableui

import (
	"github.com/aldernero/timebox/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

type mode int

const (
	boxSummary mode = iota
	boxView
	timeline
)

type Model struct {
	mode    mode
	boxName string
	tp      util.Period
	tb      *util.TimeBox
	table   table.Model
}

func New(tb *util.TimeBox, period util.Period) Model {
	return Model{mode: boxSummary, tp: period, tb: tb, table: makeBoxSummaryTable(tb, period)}
}

func (m Model) SetBoxSummary() {
	m.mode = boxSummary
	m.boxName = ""
}

func (m Model) SetBoxView(boxName string) {
	m.mode = boxView
	m.boxName = boxName
}

func (m Model) SetTimeline() {
	m.mode = timeline
	m.boxName = ""
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg.(type) {
	case modeChangedMsg:
		switch m.mode {
		case boxSummary:
			m.table = makeBoxSummaryTable(m.tb, m.tp)
		case boxView:
			m.table = makeBoxViewTable(m.tb, m.tp, m.boxName)
		case timeline:
			m.table = makeTimelineTable(m.tb, m.tp)
		}
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}
