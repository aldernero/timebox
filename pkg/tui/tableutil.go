package tui

import (
	util2 "github.com/aldernero/timebox/pkg/util"
	"github.com/evertras/bubble-table/table"
	"time"
)

const (
	columnKeyBox    = "box"
	columnKeyMin    = "min"
	columnKeyMax    = "max"
	columnKeyUse    = "use"
	columnKeyStart  = "start"
	columnKeyEnd    = "end"
	columnKeyDur    = "dur"
	columnWidthBox  = 24
	columnWidthTime = 20
	columnWidthDur  = 12
)

func makeBoxSummaryRow(box string, min, max, use time.Duration) table.Row {
	return table.NewRow(table.RowData{
		columnKeyBox: box,
		columnKeyMin: min,
		columnKeyMax: max,
		columnKeyUse: util2.DurationParser(use),
	})
}

func makeTimelineRow(box string, start time.Time, end time.Time) table.Row {
	return table.NewRow(table.RowData{
		columnKeyBox:   box,
		columnKeyStart: start.Format(time.DateTime),
		columnKeyEnd:   end.Format(time.DateTime),
		columnKeyDur:   util2.DurationParser(end.Sub(start)),
	})
}

func makeBoxSummaryTable(tb util2.TimeBox, p util2.Period) table.Model {
	boxes := tb.Boxes
	var rows []table.Row
	timespan := util2.PeriodSoFar(p, time.January)
	for _, val := range tb.Names {
		box := boxes[val]
		minTime, maxTime := box.ScaledTimes(p)
		spans := tb.GetSpansForBox(val, timespan)
		usedTime := spans.Duration()
		rows = append(rows, makeBoxSummaryRow(val, minTime, maxTime, usedTime))
	}
	return table.New([]table.Column{
		table.NewFlexColumn(columnKeyBox, "Box", 2),
		table.NewFlexColumn(columnKeyMin, "Min", 1),
		table.NewFlexColumn(columnKeyMax, "Max", 1),
		table.NewFlexColumn(columnKeyUse, "Used", 1),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(TableStyle).
		WithTargetWidth(UIWidth).
		WithPageSize(SummaryPageSize).
		Focused(true)
}

func makeBoxViewTable(tb util2.TimeBox, boxName string, p util2.Period) table.Model {
	var rows []table.Row
	timespan := util2.PeriodSoFar(p, time.January)
	spans := tb.GetSpansForBox(boxName, timespan)
	for _, val := range spans.Spans {
		rows = append(rows, makeTimelineRow(boxName, val.Start, val.End))
	}
	return table.New([]table.Column{
		table.NewFlexColumn(columnKeyBox, "Box", 2),
		table.NewColumn(columnKeyStart, "Start", columnWidthTime),
		table.NewColumn(columnKeyEnd, "End", columnWidthTime),
		table.NewFlexColumn(columnKeyDur, "Duration", 1),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(TableStyle).
		WithTargetWidth(UIWidth).
		WithPageSize(SummaryPageSize).
		Focused(true)
}

func makeTimelineTable(tb util2.TimeBox, p util2.Period) table.Model {
	var rows []table.Row
	timespan := util2.PeriodSoFar(p, time.January)
	spans := tb.GetSpansForTimespan(timespan)
	for _, val := range spans.Spans {
		rows = append(rows, makeTimelineRow(val.Box, val.Start, val.End))
	}
	return table.New([]table.Column{
		table.NewFlexColumn(columnKeyBox, "Box", columnWidthBox),
		table.NewFlexColumn(columnKeyStart, "Start", columnWidthDur),
		table.NewFlexColumn(columnKeyEnd, "End", columnWidthDur),
		table.NewFlexColumn(columnKeyDur, "Duration", columnWidthDur),
	}).WithRows(rows).
		BorderRounded().
		WithBaseStyle(TableStyle).
		WithTargetWidth(UIWidth).
		WithPageSize(SummaryPageSize).
		Focused(true)
}
