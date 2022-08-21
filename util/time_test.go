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

func TestMonthStart(t *testing.T) {
	tests := map[string]struct {
		year  int
		month int
		day   int
		hour  int
		min   int
		sec   int
	}{
		"year start":  {2022, 1, 5, 0, 0, 0},
		"month start": {2022, 2, 1, 0, 0, 0},
		"year end":    {2021, 12, 31, 23, 59, 59},
		"month end":   {2021, 5, 31, 20, 45, 0},
		"middle":      {1997, 7, 10, 8, 17, 51},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t1 := time.Date(tc.year, time.Month(tc.month), tc.day, tc.hour, tc.min, tc.sec, 0, time.Local)
			t2 := MonthStart(t1)
			assert.Equal(t, 1, t2.Day())
			assert.Zero(t, t2.Hour())
			assert.Zero(t, t2.Minute())
			assert.Zero(t, t2.Second())
		})
	}
}

func TestQuarterStart(t *testing.T) {
	tests := map[string]struct {
		year     int
		month    int
		day      int
		hour     int
		min      int
		sec      int
		fys      time.Month
		expMonth time.Month
	}{
		"year start":  {2022, 1, 5, 0, 0, 0, time.January, time.January},
		"month start": {2022, 2, 1, 0, 0, 0, time.March, time.December},
		"year end":    {2021, 12, 31, 23, 59, 59, time.June, time.December},
		"month end":   {2021, 5, 31, 20, 45, 0, time.September, time.March},
		"middle":      {1997, 7, 10, 8, 17, 51, time.November, time.May},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t1 := time.Date(tc.year, time.Month(tc.month), tc.day, tc.hour, tc.min, tc.sec, 0, time.Local)
			t2 := QuarterStart(t1, tc.fys)
			assert.Equal(t, tc.expMonth, t2.Month())
			assert.Equal(t, 1, t2.Day())
			assert.Zero(t, t2.Hour())
			assert.Zero(t, t2.Minute())
			assert.Zero(t, t2.Second())
		})
	}
}

func TestFiscalQuarter(t *testing.T) {
	tests := map[string]struct {
		fiscalYearStart time.Month
		monthInQuestion time.Month
		expectedQuarter int
	}{
		"fys1m1q1":   {time.January, time.January, 1},
		"fys1m2q1":   {time.January, time.February, 1},
		"fys1m3q1":   {time.January, time.March, 1},
		"fys1m4q2":   {time.January, time.April, 2},
		"fys1m7q3":   {time.January, time.July, 3},
		"fys1m11q4":  {time.January, time.December, 4},
		"fys2m1q1":   {time.February, time.January, 4},
		"fys2m2q1":   {time.February, time.February, 1},
		"fys2m7q2":   {time.February, time.July, 2},
		"fys4m10q3":  {time.April, time.October, 3},
		"fys8m5q4":   {time.August, time.May, 4},
		"fys10m12q1": {time.October, time.December, 1},
		"fys12m12q1": {time.December, time.December, 1},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedQuarter, FiscalQuarter(tc.fiscalYearStart, tc.monthInQuestion))
		})
	}
}
