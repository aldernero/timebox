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

func TimeBoxFromDB(dbname string) TimeBox {
	var tb TimeBox
	tb.Fname = dbname
	tb.tbdb = db.NewDBWithName(dbname)
	tb.tbdb.Init()
	tb.Names, tb.Boxes = AllBoxesFromDB(tb.tbdb)
	tb.Spans = AllSpansFromDB(tb.tbdb)
	return tb
}

func (tb TimeBox) GetSpansForBox(box string, span Span) SpanSet {
	spans := NewSpanSet()
	boxSpans := tb.Spans[box]

	for _, s := range boxSpans.Spans {
		overlap := s.GetOverlap(span)
		if !overlap.IsZero() {
			spans.Add(overlap)
		}
	}
	return spans
}

func (tb TimeBox) GetSpansForTimespan(span Span) SpanSet {
	spans := NewSpanSet()
	start := span.Start.Unix()
	end := span.End.Unix()
	spanRow, err := tb.tbdb.GetSpansForTimeRange(start, end)
	if err != nil {
		panic(err)
	}
	for _, sr := range spanRow {
		spans.Add(Span{time.Unix(sr.Start, 0), time.Unix(sr.End, 0), sr.Box})
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
				spanset := spans[box]
				spanset.Add(overlap)
				spans[box] = spanset
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
		tb.Spans[box] = NewSpanSet()
	}
	spanset := tb.Spans[box]
	spanset.Add(span)
	tb.Spans[box] = spanset
	return nil
}

func (tb TimeBox) DeleteSpan(span Span) error {
	box := span.Box
	err := tb.tbdb.DeleteSpan(span.Start.Unix(), span.End.Unix(), box)
	if err != nil {
		return err
	}
	spanset := tb.Spans[box]
	spanset.Remove(span)
	tb.Spans[box] = spanset
	return nil
}
