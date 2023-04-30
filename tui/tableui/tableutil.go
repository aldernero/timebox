package tableui

import (
	"github.com/aldernero/timebox/tui/constants"
	"github.com/aldernero/timebox/util"
	"github.com/evertras/bubble-table/table"
	"time"
)

const (
	columnKeyBox   = "box"
	columnKeyMin   = "min"
	columnKeyMax   = "max"
	columnKeyUse   = "use"
	columnWidthBox = 20
	columnWidthDur = 10
)

func makeRow(box string, min, max, use time.Duration) table.Row {
	return table.NewRow(table.RowData{
		columnKeyBox: box,
		columnKeyMin: min,
		columnKeyMax: max,
		columnKeyUse: use,
	})
}

func makeTable(tb *util.TimeBox, p util.Period) table.Model {
	boxes := tb.Boxes
	var rows []table.Row
	timespan := util.PeriodSoFar(p, time.January)
	for _, val := range tb.Names {
		box := boxes[val]
		minTime, maxTime := box.ScaledTimes(p)
		usedTime := tb.GetSpansForBox(val, timespan).Duration()
		rows = append(rows, makeRow(val, minTime, maxTime, usedTime))
	}
	return table.New([]table.Column{
		table.NewFlexColumn(columnKeyBox, "Box", columnWidthBox),
		table.NewFlexColumn(columnKeyMin, "Min", columnWidthDur),
		table.NewFlexColumn(columnKeyMax, "Max", columnWidthDur),
		table.NewFlexColumn(columnKeyUse, "Used", columnWidthDur),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(constants.TableStyle).
		WithTargetWidth(constants.TUIWidth).
		WithPageSize(constants.SummaryPageSize).
		Focused(true)
}
