package tableui

import (
	"github.com/aldernero/timebox/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"time"
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
	tbl     table.Model
}

func NewBoxSummary(tb *util.TimeBox, period util.Period) Model {
	return Model{mode: boxSummary, tp: period, tb: tb, tbl: makeBoxSummaryTable(tb, period)}
}

func NewBoxView(tb *util.TimeBox, period util.Period, boxName string) Model {
	return Model{mode: boxView, boxName: boxName, tp: period, tb: tb, tbl: makeBoxViewTable(tb, period, boxName)}
}

func NewTimeline(tb *util.TimeBox, period util.Period) Model {
	return Model{mode: timeline, tp: period, tb: tb, tbl: makeTimelineTable(tb, period)}
}

func (m Model) GetSelectedBoxName() string {
	row := m.tbl.HighlightedRow()
	data := row.Data[columnKeyBox].(string)
	return data
}

func (m Model) GetSelectedSpan() util.Span {
	row := m.tbl.HighlightedRow()
	box := row.Data[columnKeyBox].(string)
	start := row.Data[columnKeyStart].(time.Time)
	end := row.Data[columnKeyEnd].(time.Time)
	return util.Span{Start: start, End: end, Name: box}
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
			m.tbl = makeBoxSummaryTable(m.tb, m.tp)
		case boxView:
			m.tbl = makeBoxViewTable(m.tb, m.tp, m.boxName)
		case timeline:
			m.tbl = makeTimelineTable(m.tb, m.tp)
		}
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.tbl.View()
}
