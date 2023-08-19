package util

import (
	"github.com/aldernero/timebox/db"
	"time"
)

type Span struct {
	Start time.Time
	End   time.Time
}

func (s Span) Duration() time.Duration {
	return s.End.Sub(s.Start)
}

func (s Span) IsZero() bool {
	return s.Start.IsZero() && s.End.IsZero()
}

func (s Span) Overlaps(span Span) bool {
	disjoint := s.Start.After(span.End) || s.End.Before(s.Start)
	return !disjoint
}

func (s Span) GetOverlap(span Span) Span {
	var result Span
	if !s.Overlaps(span) {
		result.Start = Later(s.Start, span.Start)
		result.End = Earlier(s.End, span.End)
	}
	return result
}

type SpanSet struct {
	Spans []Span
}

func (s SpanSet) IsEmpty() bool {
	return len(s.Spans) == 0
}

func (s SpanSet) Add(span Span) {
	s.Spans = append(s.Spans, span)
}

func (s SpanSet) Duration() time.Duration {
	var seconds int64
	for _, span := range s.Spans {
		seconds += span.End.Unix() - span.Start.Unix()
	}
	return time.Duration(seconds) * time.Second
}

func AllSpansFromDB(tbdb db.TBDB) map[string]SpanSet {
	result := make(map[string]SpanSet)
	brs, err := tbdb.GetAllBoxes()
	if err != nil {
		panic(err)
	}
	for _, br := range brs {
		spanset := SpanSet{}
		srs, err := tbdb.GetSpansForBox(br.Name)
		if err != nil {
			panic(err)
		}
		for _, sr := range srs {
			spanset.Add(Span{
				Start: time.Unix(sr.Start, 0),
				End:   time.Unix(sr.End, 0),
			})
		}
		result[br.Name] = spanset
	}
	return result
}

func AllSpansFromDBForTimeRange(tbdb db.TBDB, start, end int64) SpanSet {
	result := SpanSet{}
	srs, err := tbdb.GetSpansForTimeRange(start, end)
	if err != nil {
		panic(err)
	}
	for _, sr := range srs {
		result.Add(Span{
			Start: time.Unix(sr.Start, 0),
			End:   time.Unix(sr.End, 0),
		})
	}
	return result
}
