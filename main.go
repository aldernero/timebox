package main

import (
	"github.com/aldernero/timebox/tui"
	"github.com/aldernero/timebox/util"
)

var dbName = "timebox.db"

func main() {

	timebox := util.TimeBoxFromDB(dbName)
	tui.StartTea(timebox)
}
