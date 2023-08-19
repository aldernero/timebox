package util

import (
	"fmt"
	"time"
)

type Period int

const (
	Week Period = iota
	Month
	Quarter
	Year
)

const (
	secondsPerYear   = 31557600 // 365.25 days
	secondsPerWeek   = 604800
	secondsPerDay    = 86400
	secondsPerHour   = 3600
	secondsPerMinute = 60
	monthsPerWeek    = 0.2301 // estimate for faster calculation
	quartersPerWeek  = 0.0767 // estimate for faster calculation
	yearsPerWeek     = 0.0192 // estimate for faster calculation
)

type TimePeriod struct {
	Period
}

func (t *TimePeriod) Names() []string {
	names := []string{"Week", "Month", "Quarter", "Year"}
	return names
}

func (t *TimePeriod) Current() int {
	return int(t.Period)
}

func (t *TimePeriod) Next() {
	t.Period = (t.Period + 1) % 4
}

func (t *TimePeriod) Previous() {
	t.Period = (t.Period + 3) % 4
}

// WeekStart calculates the time at the beginning of the week for a given time
//
//goland:noinspection SpellCheckingInspection
func WeekStart(t time.Time) time.Time {
	wday := int(t.Weekday())
	seconds := wday*86400 + t.Hour()*3600 + t.Minute()*60 + t.Second()
	return t.Add(time.Duration(-seconds) * time.Second)
}

// ThisWeekStart calculates the time at the beginning of the current week
func ThisWeekStart() time.Time {
	return WeekStart(time.Now())
}

// MonthStart calculates the time at the beginning of the month for a given time
func MonthStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

// ThisMonthStart calculates the time at the beginning of the current month
func ThisMonthStart() time.Time {
	return MonthStart(time.Now())
}

func QuarterStart(t time.Time, fys time.Month) time.Time {
	y := t.Year()
	m := int(t.Month())
	diff := m - int(fys)
	if diff < 0 {
		diff *= -1
	}
	monthOffset := (3 - diff%3) % 3
	m -= monthOffset
	if m < 1 {
		y--
		m += 12
	}
	return time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.Local)
}

func ThisQuarterStart(fys time.Month) time.Time {
	return QuarterStart(time.Now(), fys)
}

// YearStart calculates the time at the beginning of the month for a given time
func YearStart(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
}

// ThisYearStart calculates the time at the beginning of the current year
func ThisYearStart() time.Time {
	return YearStart(time.Now())
}

//goland:noinspection SpellCheckingInspection
func DurationParser(d time.Duration) string {
	dsec := int(d.Seconds())
	years := dsec / secondsPerYear
	days := (dsec - years*secondsPerYear) / secondsPerDay
	hours := (dsec - years*secondsPerYear - days*secondsPerDay) / secondsPerHour
	minutes := (dsec - years*secondsPerYear - days*secondsPerDay - hours*secondsPerHour) / secondsPerMinute
	seconds := dsec - years*secondsPerYear - days*secondsPerDay - hours*secondsPerHour - minutes*secondsPerMinute
	var result string
	if years > 0 {
		result = fmt.Sprintf("%dy %dd %dh %dm %ds", years, days, hours, minutes, seconds)
	}
	if years == 0 && days > 0 {
		result = fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if years == 0 && days == 0 && hours > 0 {
		result = fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if years == 0 && days == 0 && hours == 0 && minutes > 0 {
		result = fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	if years == 0 && days == 0 && hours == 0 && minutes == 0 {
		result = fmt.Sprintf("%ds", seconds)
	}
	return result
}

func Earlier(t1, t2 time.Time) time.Time {
	if t1.Unix() < t2.Unix() {
		return t1
	}
	return t2
}

func Later(t1, t2 time.Time) time.Time {
	if t1.Unix() > t2.Unix() {
		return t1
	}
	return t2
}

func WeekSoFar() Span {
	return Span{Start: ThisWeekStart(), End: time.Now()}
}

func MonthSoFar() Span {
	return Span{Start: ThisMonthStart(), End: time.Now()}
}

func QuarterSoFar(fys time.Month) Span {
	return Span{Start: ThisQuarterStart(fys), End: time.Now()}
}

func YearSoFar() Span {
	return Span{Start: ThisYearStart(), End: time.Now()}
}

func PeriodSoFar(p Period, fys time.Month) Span {
	var result Span
	switch p {
	case Week:
		result = WeekSoFar()
	case Month:
		result = MonthSoFar()
	case Quarter:
		result = QuarterSoFar(fys)
	case Year:
		result = YearSoFar()
	}
	return result
}

func FiscalQuarter(fiscalYearStart, calendarMonth time.Month) int {
	fm := int(calendarMonth - fiscalYearStart)
	if fm < 0 {
		fm += 12
	}
	return fm/3 + 1
}

func AfterOrEqual(t1, t2 time.Time) bool {
	return t1.Unix() >= t2.Unix()
}

func BeforeOrEqual(t1, t2 time.Time) bool {
	return t1.Unix() <= t2.Unix()
}
