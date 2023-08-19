package util

import (
	"github.com/aldernero/timebox/db"
	"time"
)

type TimeBox struct {
	tbdb  db.TBDB
	Fname string
	Names []string
	Boxes map[string]Box
	Spans map[string]SpanSet
}

func NewTimeBox(dbname string) TimeBox {
	var tb TimeBox
	tb.Fname = dbname
	tb.tbdb = db.NewDBWithName(dbname)
	tb.tbdb.Init()
	tb.Names, tb.Boxes = AllBoxesFromDB(tb.tbdb)
	if len(tb.Names) == 0 {
		panic("no boxes")
	}
	tb.Spans = AllSpansFromDB(tb.tbdb)
	if len(tb.Spans) == 0 {
		panic("no spans")
	}
	return tb
}

func (tb TimeBox) GetSpansForBox(box string, span Span) SpanSet {
	var spans SpanSet
	boxSpans := tb.Spans[box]

	for _, s := range boxSpans.Spans {
		overlap := s.GetOverlap(span)
		if !overlap.IsZero() {
			spans.Add(overlap)
		}
		spans.Add(span)
	}
	return spans
}

func (tb TimeBox) GetSpansForTimespan(span Span) SpanSet {
	var spans SpanSet
	start := span.Start.Unix()
	end := span.End.Unix()
	spanRow, err := tb.tbdb.GetSpansForTimeRange(start, end)
	if err != nil {
		panic(err)
	}
	for _, sr := range spanRow {
		spans.Add(Span{time.Unix(sr.Start, 0), time.Unix(sr.End, 0)})
	}
	return spans
}

func (tb TimeBox) GetSpans(span Span) map[string]SpanSet {
	spans := make(map[string]SpanSet)
	for box, spanset := range tb.Spans {
		spans[box] = SpanSet{}
		for _, s := range spanset.Spans {
			overlap := s.GetOverlap(span)
			if !overlap.IsZero() {
				spans[box].Add(overlap)
			}
		}
	}
	return spans
}

func (tb TimeBox) AddBox(box Box) error {
	err := tb.tbdb.AddBox(box.Name, int64(box.MinTime.Seconds()), int64(box.MaxTime.Seconds()))
	if err != nil {
		return err
	}
	tb.Boxes[box.Name] = box
	return nil
}

func (tb TimeBox) UpdateBox(box Box) error {
	err := tb.tbdb.UpdateBox(box.Name, int64(box.MinTime.Seconds()), int64(box.MaxTime.Seconds()))
	if err != nil {
		return err
	}
	tb.Boxes[box.Name] = box
	return nil
}

func (tb TimeBox) DeleteBox(box string) error {
	err := tb.tbdb.DeleteBox(box)
	if err != nil {
		return err
	}
	delete(tb.Boxes, box)
	return nil
}

func (tb TimeBox) AddSpan(span Span, box string) error {
	err := tb.tbdb.AddSpan(span.Start.Unix(), span.End.Unix(), box)
	if err != nil {
		return err
	}
	if _, ok := tb.Spans[box]; !ok {
		tb.Spans[box] = SpanSet{}
	}
	tb.Spans[box].Add(span)
	return nil
}
