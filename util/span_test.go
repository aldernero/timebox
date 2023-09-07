package util

import (
	"fmt"
	"github.com/aldernero/timebox/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const dbName = "test.db"

func setup(t *testing.T) db.TBDB {
	tempDir, err := os.MkdirTemp(os.TempDir(), "timebox")
	require.NoError(t, err)
	testdb := filepath.Join(tempDir, filepath.FromSlash(dbName))
	err = os.Remove(testdb)
	require.NoFileExists(t, testdb)
	tbdb := db.NewDBWithName(testdb)
	require.NoError(t, tbdb.CreateDB())
	return tbdb
}

func TestSpan_Duration(t *testing.T) {
	tests := map[string]struct {
		span Span
		want time.Duration
	}{
		"zero": {
			span: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
			},
			want: 0,
		},
		"one second": {
			span: Span{
				Start: time.Date(2023, time.February, 14, 14, 32, 54, 0, time.Local),
				End:   time.Date(2023, time.February, 14, 14, 32, 55, 0, time.Local),
			},
			want: time.Second,
		},
		"one minute": {
			span: Span{
				Start: time.Date(2023, time.March, 3, 17, 12, 07, 0, time.Local),
				End:   time.Date(2023, time.March, 3, 17, 13, 07, 0, time.Local),
			},
			want: time.Minute,
		},
		"one hour": {
			span: Span{
				Start: time.Date(2023, time.April, 5, 23, 59, 59, 0, time.Local),
				End:   time.Date(2023, time.April, 6, 0, 59, 59, 0, time.Local),
			},
			want: time.Hour,
		},
		"one day": {
			span: Span{
				Start: time.Date(2023, time.May, 7, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.May, 8, 0, 0, 0, 0, time.Local),
			},
			want: 24 * time.Hour,
		},
		"one week": {
			span: Span{
				Start: time.Date(2023, time.June, 9, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.June, 16, 0, 0, 0, 0, time.Local),
			},
			want: 7 * 24 * time.Hour,
		},
		"one month": {
			span: Span{
				Start: time.Date(2023, time.July, 11, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.August, 11, 0, 0, 0, 0, time.Local),
			},
			want: 31 * 24 * time.Hour,
		},
		"one year": {
			span: Span{
				Start: time.Date(2023, time.September, 13, 0, 0, 0, 0, time.Local),
				End:   time.Date(2024, time.September, 13, 0, 0, 0, 0, time.Local),
			},
			want: 366 * 24 * time.Hour,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.span.Duration()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSpan_IsZero(t *testing.T) {
	tests := map[string]struct {
		span Span
		want bool
	}{
		"zero": {
			span: Span{},
			want: true,
		},
		"non-zero one": {
			span: Span{
				Start: time.Unix(0, 0),
				End:   time.Date(2023, time.January, 1, 0, 0, 1, 0, time.Local),
			},
			want: false,
		},
		"non-zero both": {
			span: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 1, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 0, 0, 2, 0, time.Local),
			},
			want: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.span.IsZero()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSpan_Overlaps(t *testing.T) {
	// Let spans be defined as
	// span1 = [s1, e1]
	// span2 = [s2, e2]
	tests := map[string]struct {
		span1 Span
		span2 Span
		want  bool
	}{
		"s2 == e1": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
			},
			want: false,
		},
		"s1 == e2": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
			},
			want: false,
		},
		"inside": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 3, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
			},
			want: true,
		},
		"overlap right": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 3, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 4, 0, 0, 0, time.Local),
			},
			want: true,
		},
		"overlap left": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 4, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 3, 0, 0, 0, time.Local),
			},
			want: true,
		},
		"overlap both": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 4, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 5, 0, 0, 0, time.Local),
			},
			want: true,
		},
		"identical spans": {
			span1: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 5, 0, 0, 0, time.Local),
			},
			span2: Span{
				Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 5, 0, 0, 0, time.Local),
			},
			want: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.span1.Overlaps(tc.span2)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSpan_GetOverlap(t *testing.T) {
	morning := Span{
		Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 11, 59, 59, 0, time.Local),
	}
	evening := Span{
		Start: time.Date(2023, time.January, 1, 17, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 23, 59, 59, 0, time.Local),
	}
	breakfast := Span{
		Start: time.Date(2023, time.January, 1, 7, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 10, 0, 0, 0, time.Local),
	}
	lunch := Span{
		Start: time.Date(2023, time.January, 1, 11, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 13, 0, 0, 0, time.Local),
	}
	dinner := Span{
		Start: time.Date(2023, time.January, 1, 18, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 20, 0, 0, 0, time.Local),
	}
	tests := map[string]struct {
		span1 Span
		span2 Span
		want  Span
	}{
		"no overlap": {
			span1: morning,
			span2: evening,
			want:  Span{},
		},
		"morning + lunch": {
			span1: morning,
			span2: lunch,
			want: Span{
				Start: time.Date(2023, time.January, 1, 11, 0, 0, 0, time.Local),
				End:   time.Date(2023, time.January, 1, 11, 59, 59, 0, time.Local),
			},
		},
		"morning + breakfast": {
			span1: morning,
			span2: breakfast,
			want:  breakfast,
		},
		"evening + dinner": {
			span1: evening,
			span2: dinner,
			want:  dinner,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.span1.GetOverlap(tc.span2)
			assert.True(t, tc.want.IsEqual(got))
		})
	}
}

func TestSpanSet_Duration(t *testing.T) {
	tests := map[string]struct {
		spans SpanSet
		want  time.Duration
	}{
		"empty": {
			spans: NewSpanSet(),
			want:  0,
		},
		"one": {
			spans: SpanSet{
				Spans: []Span{
					{
						Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
						End:   time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
					},
				},
			},
			want: time.Hour,
		},
		"two": {
			spans: SpanSet{
				Spans: []Span{
					{
						Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
						End:   time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
					},
					{
						Start: time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
						End:   time.Date(2023, time.January, 1, 3, 0, 0, 0, time.Local),
					},
				},
			},
			want: 2 * time.Hour,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.spans.Duration()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestAllSpansFromDB(t *testing.T) {
	tbdb := setup(t)
	for i := 0; i < 4; i++ {
		boxName := fmt.Sprintf("box%d", i)
		err := tbdb.AddBox(boxName, 100, 1000)
		require.NoError(t, err)
		for j := 0; j < 4; j++ {
			start := time.Date(2023, time.January, 1, i, j, 0, 0, time.Local)
			end := time.Date(2023, time.January, 1, i, j+1, 0, 0, time.Local)
			err := tbdb.AddSpan(start.Unix(), end.Unix(), boxName)
			require.NoError(t, err)
		}
	}
	spans := AllSpansFromDB(tbdb)
	assert.Equal(t, 4, len(spans))
	for i := 0; i < 4; i++ {
		boxName := fmt.Sprintf("box%d", i)
		spanset := spans[boxName]
		assert.Equal(t, 4, spanset.Size())
	}
}

func TestAllSpansFromDBForTimeRange(t *testing.T) {
	tbdb := setup(t)
	span1 := Span{
		Start: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 1, 0, 0, 0, time.Local),
	}
	span2 := Span{
		Start: time.Date(2023, time.January, 1, 2, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 3, 0, 0, 0, time.Local),
	}
	span3 := Span{
		Start: time.Date(2023, time.January, 1, 4, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 5, 0, 0, 0, time.Local),
	}
	span4 := Span{
		Start: time.Date(2023, time.January, 1, 6, 0, 0, 0, time.Local),
		End:   time.Date(2023, time.January, 1, 7, 0, 0, 0, time.Local),
	}
	for i := 0; i < 4; i++ {
		boxName := fmt.Sprintf("box%d", i)
		err := tbdb.AddBox(boxName, 100, 1000)
		require.NoError(t, err)
	}
	err := tbdb.AddSpan(span1.Start.Unix(), span1.End.Unix(), "box1")
	require.NoError(t, err)
	err = tbdb.AddSpan(span2.Start.Unix(), span2.End.Unix(), "box1")
	require.NoError(t, err)
	err = tbdb.AddSpan(span3.Start.Unix(), span3.End.Unix(), "box2")
	require.NoError(t, err)
	err = tbdb.AddSpan(span4.Start.Unix(), span4.End.Unix(), "box3")
	require.NoError(t, err)
	spanRow, err := tbdb.GetSpansForTimeRange(span1.Start.Unix(), span4.End.Unix())
	require.NoError(t, err)
	assert.Equal(t, 4, len(spanRow))
	spans := AllSpansFromDBForTimeRange(tbdb, span1.Start, span4.End)
	assert.Equal(t, 4, spans.Size())
}
