package main

import (
	"github.com/aldernero/timebox/pkg/tui"
	"github.com/aldernero/timebox/pkg/util"
)

var dbName = "timebox.db"

func main() {
	timebox := util.TimeBoxFromDB(dbName)
	tui.StartTea(timebox)
}
