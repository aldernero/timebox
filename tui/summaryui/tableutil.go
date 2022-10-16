package summaryui

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
	columnWidthBox = 15
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

func makeTable(tb *util.TimeBox) table.Model {
	boxes := tb.Boxes
	var rows []table.Row
	for key, val := range boxes {
		rows = append(rows, makeRow(key, val.MinTime,
			val.MaxTime, val.MinTime))
	}
	return table.New([]table.Column{
		table.NewColumn(columnKeyBox, "Box", columnWidthBox),
		table.NewColumn(columnKeyMin, "Min", columnWidthDur),
		table.NewColumn(columnKeyMax, "Max", columnWidthDur),
		table.NewColumn(columnKeyUse, "Used", columnWidthDur),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(constants.TableStyle).
		WithPageSize(constants.SummaryPageSize).
		Focused(true)
}
