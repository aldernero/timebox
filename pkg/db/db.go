package db

import (
	"database/sql"
	"fmt"
	//_ "github.com/mattn/go-sqlite3"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

const defaultDriver = "sqlite"

// TBDB is the base struct for the database
type TBDB struct {
	name   string
	driver string
}

type SpanRow struct {
	ID    int64
	Start int64
	End   int64
	Box   string
}

type BoxRow struct {
	Name       string
	CreateTime int64
	MinTime    int64
	MaxTime    int64
}

func NewDBWithName(name string) TBDB {
	return TBDB{
		name:   name,
		driver: defaultDriver,
	}
}

func (d TBDB) Init() {
	if _, err := os.Stat(d.name); os.IsNotExist(err) {
		err := d.CreateDB()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Create functions

func (d TBDB) CreateDB() error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	// spans table, stores spans of time spent on a given box
	sqlStmt := `
	CREATE TABLE spans (id INTEGER PRIMARY KEY AUTOINCREMENT, start INTEGER NOT NULL, end INTEGER NOT NULL, box TEXT NOT NULL);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	// boxes table, stores active boxes and their configurations
	sqlStmt = `
	CREATE TABLE boxes (name TEXT NOT NULL PRIMARY KEY, createTime INTEGER NOT NULL, minTime INTEGER NOT NULL, maxTime INTEGER NOT NULL);
	`
	_, err = db.Exec(sqlStmt)
	return err
}

func (d TBDB) AddSpan(start, end int64, box string) error {
	if start > end {
		return fmt.Errorf("start time is after end time")
	}
	now := time.Now().Unix()
	if start > now || end > now {
		return fmt.Errorf("time span is in the future")
	}
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	exists, err := d.DoesBoxExist(box)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("box %s doesn't exist", box)
	}
	overlaps, err := d.DoesSpanOverlap(start, end)
	if err != nil {
		return err
	}
	if overlaps {
		return fmt.Errorf("time overlaps existing span")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO spans(start, end, box) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	_, err = stmt.Exec(start, end, box)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (d TBDB) AddBox(name string, minTime, maxTime int64) error {
	if minTime > maxTime {
		return fmt.Errorf("minTime is greater than maxTime")
	}
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO boxes(name, createTime, minTime, maxTime) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	now := time.Now().Unix()
	_, err = stmt.Exec(name, now, minTime, maxTime)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (d TBDB) DoesSpanOverlap(start, end int64) (bool, error) {
	var result bool
	var count int
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	row := db.QueryRow("SELECT COUNT(*) FROM spans WHERE NOT ? <= start AND NOT ? >= end;", end, start)
	err = row.Scan(&count)
	if err != nil {
		return result, err
	}
	result = count > 0
	return result, nil
}

func (d TBDB) DoesBoxExist(name string) (bool, error) {
	var result bool
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	row := db.QueryRow("SELECT COUNT(*) FROM boxes WHERE name = ?", name)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return result, err
	}
	result = count > 0
	return result, nil
}

// Read functions

func (d TBDB) GetBox(name string) (BoxRow, error) {
	var result BoxRow
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	row := db.QueryRow("SELECT * FROM boxes WHERE name = ?", name)
	var createTime, minTime, maxTime int64
	err = row.Scan(&name, &createTime, &minTime, &maxTime)
	if err != nil {
		return result, err
	}
	return BoxRow{name, createTime, minTime, maxTime}, nil
}

func (d TBDB) GetAllBoxes() ([]BoxRow, error) {
	var result []BoxRow
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	rows, err := db.Query("SELECT * FROM boxes ORDER BY createTime DESC")
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var name string
		var createTime, minTime, maxTime int64
		err := rows.Scan(&name, &createTime, &minTime, &maxTime)
		if err != nil {
			return result, err
		}
		result = append(result, BoxRow{name, createTime, minTime, maxTime})
	}
	return result, nil
}

func (d TBDB) GetSpansForBox(boxName string) ([]SpanRow, error) {
	var result []SpanRow
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	box, err := d.GetBox(boxName)
	if err != nil {
		return result, err
	}
	rows, err := db.Query("SELECT * FROM spans WHERE box = ? ORDER BY start", box.Name)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var id int64
		var name string
		var start, end int64
		err := rows.Scan(&id, &start, &end, &name)
		if err != nil {
			return result, err
		}
		result = append(result, SpanRow{id, start, end, name})
	}
	return result, nil
}

func (d TBDB) GetSpansForTimeRange(start, end int64) ([]SpanRow, error) {
	var result []SpanRow
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	rows, err := db.Query("SELECT * FROM spans WHERE start >= ? AND end <= ? ORDER BY start", start, end)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		id := -1
		var name string
		var start, end int64
		err := rows.Scan(&id, &start, &end, &name)
		if err != nil {
			return result, err
		}
		result = append(result, SpanRow{int64(id), start, end, name})
	}
	return result, nil
}

// Update functions

func (d TBDB) UpdateBox(name string, minTime, maxTime int64) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	stmt, err := db.Prepare("UPDATE boxes SET minTime = ?, maxTime = ? WHERE name = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(minTime, maxTime, name)
	return err
}

// Delete functions

func (d TBDB) DeleteBox(name string) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	stmt, err := db.Prepare("DELETE FROM boxes WHERE name = ?")
	_, err = stmt.Exec(name)
	return err
}

func (d TBDB) DeleteBoxAndSpans(name string) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	stmt, err := db.Prepare("DELETE FROM boxes WHERE name = ?")
	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare("DELETE FROM spans WHERE box = ?")
	_, err = stmt.Exec(name)
	return err
}

func (d TBDB) DeleteSpan(start, end int64, box string) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	stmt, err := db.Prepare("DELETE FROM spans WHERE start = ? AND end = ? AND box = ? ")
	_, err = stmt.Exec(start, end, box)
	return err
}

func (d TBDB) DeleteSpanByID(id int64) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	stmt, err := db.Prepare("DELETE FROM spans WHERE id = ? ")
	_, err = stmt.Exec(id)
	return err
}

func (d TBDB) UpdateSpan(id, start, end int64, box string) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	stmt, err := db.Prepare("UPDATE spans SET start = ?, end = ?, box = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(start, end, box, id)
	return err
}
