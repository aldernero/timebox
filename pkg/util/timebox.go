package util

import (
	"errors"
	"github.com/aldernero/timebox/pkg/db"
	"time"
)

type TimeBox struct {
	tbdb      db.TBDB
	Fname     string
	Names     []string
	Boxes     map[string]Box
	SpansSets map[string]SpanSet
	Spans     map[int64]Span
}

func TimeBoxFromDB(dbname string) TimeBox {
	var tb TimeBox
	tb.Fname = dbname
	tb.tbdb = db.NewDBWithName(dbname)
	tb.tbdb.Init()
	tb.Names, tb.Boxes = AllBoxesFromDB(tb.tbdb)
	tb.SpansSets, tb.Spans = AllSpansFromDB(tb.tbdb)
	return tb
}

func (tb TimeBox) SyncFromDB() {
	tb.tbdb = db.NewDBWithName(tb.Fname)
	tb.tbdb.Init()
	tb.Names, tb.Boxes = AllBoxesFromDB(tb.tbdb)
	tb.SpansSets, tb.Spans = AllSpansFromDB(tb.tbdb)
}

func (tb TimeBox) GetSpansForBox(box string, span Span) SpanSet {
	spans := NewSpanSet()
	boxSpans := tb.SpansSets[box]

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
		spans.Add(Span{ID: sr.ID, Start: time.Unix(sr.Start, 0), End: time.Unix(sr.End, 0), Box: sr.Box})
	}
	return spans
}

func (tb TimeBox) GetSpans(span Span) map[string]SpanSet {
	spans := make(map[string]SpanSet)
	for box, spanset := range tb.SpansSets {
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

func (tb TimeBox) DeleteBoxAndSpans(box string) error {
	err := tb.tbdb.DeleteBoxAndSpans(box)
	if err != nil {
		return err
	}
	tb.SyncFromDB()
	return nil
}

func (tb TimeBox) AddSpan(span Span, box string) error {
	err := tb.tbdb.AddSpan(span.Start.Unix(), span.End.Unix(), box)
	if err != nil {
		return err
	}
	if _, ok := tb.SpansSets[box]; !ok {
		tb.SpansSets[box] = NewSpanSet()
	}
	spanset := tb.SpansSets[box]
	spanset.Add(span)
	tb.SpansSets[box] = spanset
	return nil
}

func (tb TimeBox) DeleteSpan(span Span) error {
	box := span.Box
	err := tb.tbdb.DeleteSpan(span.Start.Unix(), span.End.Unix(), box)
	if err != nil {
		return err
	}
	spanset := tb.SpansSets[box]
	spanset.Remove(span)
	tb.SpansSets[box] = spanset
	return nil
}

func (tb TimeBox) DeleteSpanByID(id int64) error {
	err := tb.tbdb.DeleteSpanByID(id)
	if err != nil {
		return err
	}
	tb.SyncFromDB()
	return nil
}

func (tb TimeBox) UpdateSpan(span Span) error {
	// check if span overlaps with any other spans
	for k, v := range tb.Spans {
		if k == span.ID {
			continue
		}
		if v.Overlaps(span) {
			return errors.New("updated span overlaps with an existing span")
		}
	}
	err := tb.tbdb.UpdateSpan(span.ID, span.Start.Unix(), span.End.Unix(), span.Box)
	if err != nil {
		return err
	}
	tb.SyncFromDB()
	return nil
}
