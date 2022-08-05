package main

import (
	"github.com/aldernero/timebox/db"
)

const dbName = "timebox.db"

func main() {
	tdb := db.NewDBWithName(dbName)
	tdb.Init()
}
