package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"time"
)

const defaultDriver = "sqlite3"

// TimeBox DB
type TBDB struct {
	name   string
	driver string
}

type SpanRow struct {
	start int64
	end   int64
	name  string
}

type BoxRow struct {
	name       string
	createTime int64
	minTime    int64
	maxTime    int64
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
	defer db.Close()
	// ledger table, stores spans of time spent on a given box
	sqlStmt := `
	CREATE TABLE spans (start INTEGER NOT NULL, end INTEGER NOT NULL, name TEXT NOT NULL);
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

func (d TBDB) AddSpan(start, end int64, name string) error {
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
	defer db.Close()
	exists, err := d.DoesBoxExist(name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("box %s doesn't exist", name)
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
	stmt, err := tx.Prepare("INSERT INTO spans(start, end, name) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(start, end, name)
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
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO boxes(name, createTime, minTime, maxTime) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
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
	defer db.Close()
	row := db.QueryRow("SELECT COUNT(*) FROM spans WHERE (start <= ? AND end >= ?) OR (start <= ? AND end >= ?)", start, start, end, end)
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
	defer db.Close()
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
	defer db.Close()
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
	defer db.Close()
	rows, err := db.Query("SELECT * FROM boxes")
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

func (d TBDB) GetSpansForBox(name string) ([]SpanRow, error) {
	var result []SpanRow
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return result, err
	}
	defer db.Close()
	box, err := d.GetBox(name)
	if err != nil {
		return result, err
	}
	rows, err := db.Query("SELECT * FROM spans WHERE name = ? AND start >= ?", box.name, box.createTime)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var name string
		var start, end int64
		err := rows.Scan(&start, &end, &name)
		if err != nil {
			return result, err
		}
		result = append(result, SpanRow{start, end, name})
	}
	return result, nil
}

// Update functions

func (d TBDB) UpdateBox(name string, minTime, maxTime int64) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer db.Close()
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
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM boxes WHERE name = ?")
	_, err = stmt.Exec(name)
	return err
}

func (d TBDB) DeleteBoxAndSpans(name string) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM boxes WHERE name = ?")
	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}
	stmt, err = db.Prepare("DELETE FROM spans WHERE name = ?")
	_, err = stmt.Exec(name)
	return err
}

func (d TBDB) DeleteSpan(start, end int64, name string) error {
	db, err := sql.Open(d.driver, d.name)
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM boxes WHERE start = ? AND end = ? AND name = ? ")
	_, err = stmt.Exec(start, end, name)
	return err
}
