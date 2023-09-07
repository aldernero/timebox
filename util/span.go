package util

import (
	"fmt"
	"github.com/aldernero/timebox/db"
	"time"
)

type Span struct {
	Start time.Time
	End   time.Time
	Box   string
}

func (s Span) Duration() time.Duration {
	return s.End.Sub(s.Start)
}

func (s Span) IsZero() bool {
	return s.Start.IsZero() && s.End.IsZero()
}

func (s Span) IsEqual(span Span) bool {
	return s.Start.Equal(span.Start) && s.End.Equal(span.End)
}

func (s Span) Overlaps(span Span) bool {
	disjoint := AfterOrEqual(s.Start, span.End) || BeforeOrEqual(s.End, span.Start)
	return !disjoint
}

func (s Span) GetOverlap(span Span) Span {
	var result Span
	if s.Overlaps(span) {
		result.Start = Later(s.Start, span.Start)
		result.End = Earlier(s.End, span.End)
	}
	return result
}

func (s Span) String() string {
	return fmt.Sprintf("%s: %s - %s", s.Box, s.Start.Format(time.RFC3339), s.End.Format(time.RFC3339))
}

func (s Span) Key() string {
	return fmt.Sprintf("%s-%d-%d", s.Box, s.Start.Unix(), s.End.Unix())
}

type SpanSet struct {
	Spans  []Span          // list for table loads
	lookup map[string]Span // map for fast lookups
}

func NewSpanSet() SpanSet {
	return SpanSet{
		Spans:  make([]Span, 0),
		lookup: make(map[string]Span),
	}
}

func (s *SpanSet) IsEmpty() bool {
	return len(s.Spans) == 0
}

func (s *SpanSet) Size() int {
	return len(s.Spans)
}

func (s *SpanSet) Add(span Span) {
	key := span.Key()
	if _, ok := s.lookup[key]; !ok {
		s.lookup[key] = span
		s.Spans = append(s.Spans, span)
		return
	} else {
		panic(fmt.Sprintf("Span already exists: %s", key))
	}
}

func (s *SpanSet) HasSpan(span Span) bool {
	key := span.Key()
	_, ok := s.lookup[key]
	return ok
}

func (s *SpanSet) Duration() time.Duration {
	var seconds int64
	for _, span := range s.Spans {
		seconds += span.End.Unix() - span.Start.Unix()
	}
	return time.Duration(seconds) * time.Second
}

func (s *SpanSet) Remove(span Span) {
	key := span.Key()
	if _, ok := s.lookup[key]; ok {
		delete(s.lookup, key)
		for i, val := range s.Spans {
			if val.Key() == key {
				s.Spans = append(s.Spans[:i], s.Spans[i+1:]...)
				return
			}
		}
	}
}

func AllSpansFromDB(tbdb db.TBDB) map[string]SpanSet {
	result := make(map[string]SpanSet)
	brs, err := tbdb.GetAllBoxes()
	if err != nil {
		panic(err)
	}
	for _, br := range brs {
		spanset := NewSpanSet()
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

func AllSpansFromDBForTimeRange(tbdb db.TBDB, start, end time.Time) SpanSet {
	result := NewSpanSet()
	srs, err := tbdb.GetSpansForTimeRange(start.Unix(), end.Unix())
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
