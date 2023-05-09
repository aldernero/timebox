package util

import "time"

const (
	timeFormatShort = "2006-01-02"
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
