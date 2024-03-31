package util

import (
	"github.com/aldernero/timebox/pkg/db"
	"time"
)

type Box struct {
	Name    string
	MinTime time.Duration
	MaxTime time.Duration
}

func AllBoxesFromDB(tbdb db.TBDB) ([]string, map[string]Box) {
	result := make(map[string]Box)
	brs, err := tbdb.GetAllBoxes()
	if err != nil {
		panic(err)
	}
	names := make([]string, len(brs))
	for i, br := range brs {
		names[i] = br.Name
		result[br.Name] = Box{
			Name:    br.Name,
			MinTime: time.Duration(br.MinTime) * time.Second,
			MaxTime: time.Duration(br.MaxTime) * time.Second,
		}
	}
	return names, result
}

func (b Box) ScaledTimes(p Period) (time.Duration, time.Duration) {
	var factor float64
	minSec := b.MinTime.Seconds()
	maxSec := b.MaxTime.Seconds()
	switch p {
	case Week:
		factor = 1.0
	case Month:
		factor = monthsPerWeek
	case Quarter:
		factor = quartersPerWeek
	case Year:
		factor = yearsPerWeek
	default:
		factor = 1.0
	}
	newMin := time.Duration(minSec/factor) * time.Second
	newMax := time.Duration(maxSec/factor) * time.Second
	return newMin, newMax
}
