package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWeekStart(t *testing.T) {
	tests := map[string]struct {
		year  int
		month int
		day   int
		hour  int
		min   int
		sec   int
	}{
		"year start":   {2022, 1, 1, 0, 0, 0},
		"month start":  {2022, 2, 1, 0, 0, 0},
		"sunday start": {2021, 12, 19, 0, 0, 0},
		"sunday end":   {2021, 7, 11, 23, 59, 59},
		"saturday end": {2021, 7, 10, 23, 59, 59},
		"middle week":  {2022, 7, 15, 13, 45, 9},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t1 := time.Date(tc.year, time.Month(tc.month), tc.day, tc.hour, tc.min, tc.sec, 0, time.Local)
			t2 := WeekStart(t1)
			assert.True(t, t2.After(t1.Add(-168*time.Hour))) // within the last week
			assert.Zero(t, t2.Hour())
			assert.Zero(t, t2.Minute())
			assert.Zero(t, t2.Second())
			assert.Equal(t, time.Sunday, t2.Weekday())
		})
	}
}
