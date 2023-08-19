package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
