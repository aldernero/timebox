package util

import (
	"fmt"
	"time"
)

type period int

const (
	Week period = iota
	Month
	Quarter
	Year
)

const (
	secondsPerYear   = 31557600
	secondsPerDay    = 86400
	secondsPerHour   = 3600
	secondsPerMinute = 60
)

func LastDays(n int) Span {
	end := time.Now()
	start := end.Add(time.Duration(-n*24) * time.Hour)
	return Span{start, end}
}

func LastWeeks(n int) Span {
	end := time.Now()
	start := end.Add(time.Duration(-n*24*7) * time.Hour)
	return Span{start, end}
}

// WeekStart calculates the time at the beginning of the week for a given time
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
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.Local)
}

// ThisMonthStart calculates the time at the beginning of the current month
func ThisMonthStart() time.Time {
	return MonthStart(time.Now())
}

// YearStart calculates the time at the beginning of the month for a given time
func YearStart(t time.Time) time.Time {
	return time.Date(t.Year(), time.January, 0, 0, 0, 0, 0, time.Local)
}

// ThisMonthStart calculates the time at the beginning of the current month
func ThisYearStart() time.Time {
	return MonthStart(time.Now())
}

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
	return t1
}

func Later(t1, t2 time.Time) time.Time {
	if t1.Unix() > t2.Unix() {
		return t1
	}
	return t1
}

func WeekSoFar() Span {
	return Span{ThisWeekStart(), time.Now()}
}

func MonthSoFar() Span {
	return Span{ThisMonthStart(), time.Now()}
}

func YearSoFar() Span {
	return Span{ThisYearStart(), time.Now()}
}
