package util

import (
	"math/rand"
	"time"
)

const (
	timeFormatShort = "15:04:05"
	timeFormatLong  = "2006-01-02 15:04:05"
)

func ParseTime(s string) (time.Time, error) {
	timeFormat := timeFormatLong
	if len(s) < len(timeFormatLong) {
		timeFormat = timeFormatShort
	}
	ts, err := time.ParseInLocation(timeFormat, s, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return ts, nil
}

type InputResult struct {
	isBox bool
	box   Box
	span  Span
}

func NewInputResultBox(box Box) InputResult {
	return InputResult{isBox: true, box: box}
}

func NewInputResultSpan(span Span) InputResult {
	return InputResult{isBox: false, span: span}
}

func (r InputResult) IsBox() bool {
	return r.isBox
}

func (r InputResult) Box() Box {
	return r.box
}

func (r InputResult) Span() Span {
	return r.span
}

type PromptState int

const (
	InUse PromptState = iota
	WasCancelled
	HasResult
)

type BoxSpec struct {
	Name        string
	Min         time.Duration
	Max         time.Duration
	SpanMin     time.Duration
	SpanMax     time.Duration
	HourMin     int
	HourMax     int
	AllowedDays []time.Weekday
}

func (b BoxSpec) IsDayAllowed(day time.Weekday) bool {
	for _, d := range b.AllowedDays {
		if d == day {
			return true
		}
	}
	return false
}

func (b BoxSpec) RandomSpan(day time.Time) Span {
	spanDuration := randomDuration(b.SpanMin, b.SpanMax)
	startTime := randomTimeWithin(day.Add(time.Hour*time.Duration(b.HourMin)), day.Add(time.Hour*time.Duration(b.HourMax)), spanDuration)
	return Span{
		Start: startTime,
		End:   startTime.Add(spanDuration),
	}
}

var TypicalBoxSpecs = map[string]BoxSpec{
	"Work": {
		Name:        "Work",
		Min:         time.Hour * 35,
		Max:         time.Hour * 50,
		SpanMin:     time.Hour * 6,
		SpanMax:     time.Hour * 9,
		HourMin:     8,
		HourMax:     18,
		AllowedDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
	},
	"Exercise": {
		Name:        "Exercise",
		Min:         time.Hour * 2,
		Max:         time.Hour * 8,
		SpanMin:     time.Minute * 30,
		SpanMax:     time.Hour * 2,
		HourMin:     5,
		HourMax:     8,
		AllowedDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
	},
	"Piano": {
		Name:        "Piano",
		Min:         time.Hour * 1,
		Max:         time.Hour * 10,
		SpanMin:     time.Minute * 30,
		SpanMax:     time.Minute * 60,
		HourMin:     17,
		HourMax:     22,
		AllowedDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
	},
	"Reading": {
		Name:        "Reading",
		Min:         time.Hour * 1,
		Max:         time.Hour * 10,
		SpanMin:     time.Minute * 15,
		SpanMax:     time.Minute * 90,
		HourMin:     17,
		HourMax:     22,
		AllowedDays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
	},
	"Hiking": {
		Name:        "Hiking",
		Min:         time.Hour * 1,
		Max:         time.Hour * 5,
		SpanMin:     time.Hour * 1,
		SpanMax:     time.Hour * 5,
		HourMin:     5,
		HourMax:     12,
		AllowedDays: []time.Weekday{time.Saturday, time.Sunday},
	},
}

func GenerateSpansFromBoxSpecs(startDay, endDay time.Time) []Span {
	var spans []Span
	day := startDay
	for day.Before(endDay) {
		for _, spec := range TypicalBoxSpecs {
			if !spec.IsDayAllowed(day.Weekday()) {
				continue
			}
			span := spec.RandomSpan(day)
			spans = append(spans, span)
		}
		day = day.Add(time.Hour * 24)
	}
	return spans
}

func randomDuration(min, max time.Duration) time.Duration {
	return min + time.Duration(rand.Int63n(int64(max-min)))
}

func randomTimeWithin(start, end time.Time, duration time.Duration) time.Time {
	startSeconds := start.Unix()
	endSeconds := end.Unix()
	spanSeconds := duration.Seconds()
	randomSeconds := rand.Int63n(endSeconds-startSeconds-int64(spanSeconds)) + startSeconds
	return time.Unix(randomSeconds, 0)
}
