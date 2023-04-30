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
	columnKeyStart = "start"
	columnKeyEnd   = "end"
	columnKeyDur   = "dur"
	columnWidthBox = 20
	columnWidthDur = 10
)

func makeBoxSummaryRow(box string, min, max, use time.Duration) table.Row {
	return table.NewRow(table.RowData{
		columnKeyBox: box,
		columnKeyMin: min,
		columnKeyMax: max,
		columnKeyUse: use,
	})
}

func makeTimelineRow(box string, start time.Time, end time.Time) table.Row {
	return table.NewRow(table.RowData{
		columnKeyBox:   box,
		columnKeyStart: start,
		columnKeyEnd:   end,
		columnKeyDur:   end.Sub(start),
	})
}

func makeBoxSummaryTable(tb *util.TimeBox, p util.Period) table.Model {
	boxes := tb.Boxes
	var rows []table.Row
	timespan := util.PeriodSoFar(p, time.January)
	for _, val := range tb.Names {
		box := boxes[val]
		minTime, maxTime := box.ScaledTimes(p)
		usedTime := tb.GetSpansForBox(val, timespan).Duration()
		rows = append(rows, makeBoxSummaryRow(val, minTime, maxTime, usedTime))
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

func makeBoxViewTable(tb *util.TimeBox, p util.Period, boxName string) table.Model {
	var rows []table.Row
	timespan := util.PeriodSoFar(p, time.January)
	spans := tb.GetSpansForBox(boxName, timespan)
	for _, val := range spans.Spans {
		rows = append(rows, makeTimelineRow(boxName, val.Start, val.End))
	}
	return table.New([]table.Column{
		table.NewFlexColumn(columnKeyBox, "Box", columnWidthBox),
		table.NewFlexColumn(columnKeyStart, "Start", columnWidthDur),
		table.NewFlexColumn(columnKeyEnd, "End", columnWidthDur),
		table.NewFlexColumn(columnKeyDur, "Duration", columnWidthDur),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(constants.TableStyle).
		WithTargetWidth(constants.TUIWidth).
		WithPageSize(constants.SummaryPageSize).
		Focused(true)
}

func makeTimelineTable(tb *util.TimeBox, p util.Period) table.Model {
	var rows []table.Row
	timespan := util.PeriodSoFar(p, time.January)
	spans := tb.GetSpansForTimespan(timespan)
	for _, val := range spans.Spans {
		rows = append(rows, makeTimelineRow(val.Name, val.Start, val.End))
	}
	return table.New([]table.Column{
		table.NewFlexColumn(columnKeyBox, "Box", columnWidthBox),
		table.NewFlexColumn(columnKeyStart, "Start", columnWidthDur),
		table.NewFlexColumn(columnKeyEnd, "End", columnWidthDur),
		table.NewFlexColumn(columnKeyDur, "Duration", columnWidthDur),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(constants.TableStyle).
		WithTargetWidth(constants.TUIWidth).
		WithPageSize(constants.SummaryPageSize).
		Focused(true)
}
