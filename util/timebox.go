package util

import (
	"github.com/aldernero/timebox/db"
	"time"
)

type TimeBox struct {
	tbdb  db.TBDB
	Names []string
	Boxes map[string]Box
	Spans map[string]SpanSet
}

func NewTimeBox(dbname string) *TimeBox {
	var tb TimeBox
	tb.tbdb = db.NewDBWithName(dbname)
	tb.tbdb.Init()
	tb.Boxes = AllBoxesFromDB(tb.tbdb)
	tb.Spans = AllSpansFromDB(tb.tbdb)
	for key, _ := range tb.Boxes {
		tb.Names = append(tb.Names, key)
	}
	return &tb
}

func (tb TimeBox) GetSpansForBox(box string, span Span) SpanSet {
	var spans SpanSet
	boxSpans := tb.Spans[box]
	for _, s := range boxSpans.Spans {
		overlap := s.GetOverlap(span)
		if !overlap.IsZero() {
			spans.Add(overlap)
		}
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

func (tb TimeBox) AddBox(name string, minTime, maxTime time.Duration) error {
	return tb.tbdb.AddBox(name, int64(minTime.Seconds()), int64(maxTime.Seconds()))
}

func (tb TimeBox) AddSpan(name string, start, end time.Time) error {
	return tb.tbdb.AddSpan(start.Unix(), end.Unix(), name)
}
