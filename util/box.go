package util

import (
	"github.com/aldernero/timebox/db"
	"time"
)

type Box struct {
	Name    string
	MinTime time.Duration
	MaxTime time.Duration
}

func AllBoxesFromDB(tbdb db.TBDB) map[string]Box {
	result := make(map[string]Box)
	brs, err := tbdb.GetAllBoxes()
	if err != nil {
		panic(err)
	}
	for _, br := range brs {
		result[br.Name] = Box{
			Name:    br.Name,
			MinTime: time.Duration(br.MinTime) * time.Second,
			MaxTime: time.Duration(br.MaxTime) * time.Second,
		}
	}
	return result
}
