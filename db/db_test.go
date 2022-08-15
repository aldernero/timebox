package db

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const dbName = "test.db"

func setup(t *testing.T) TBDB {
	tempDir, err := ioutil.TempDir(os.TempDir(), "timebox")
	require.NoError(t, err)
	testdb := filepath.Join(tempDir, filepath.FromSlash(dbName))
	os.Remove(testdb)
	tbdb := NewDBWithName(testdb)
	require.NoError(t, tbdb.CreateDB())
	return tbdb
}

func boxWithCreateTime(t *testing.T, tbdb TBDB, name string, minTime, maxTime, ts int64) {
	db, err := sql.Open(tbdb.driver, tbdb.name)
	require.NoError(t, err)
	defer db.Close()
	tx, err := db.Begin()
	require.NoError(t, err)
	stmt, err := tx.Prepare("INSERT INTO boxes(name, createTime, minTime, maxTime) values(?, ?, ?, ?)")
	require.NoError(t, err)
	defer stmt.Close()
	_, err = stmt.Exec(name, ts, minTime, maxTime)
	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)
}

func TestTBDB_CreateDB(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "timebox")
	require.NoError(t, err)
	testdb := filepath.Join(tempDir, filepath.FromSlash(dbName))
	os.Remove(testdb)
	tbdb := NewDBWithName(testdb)
	require.NoError(t, tbdb.CreateDB())
	require.FileExists(t, testdb)
	exists, err := tbdb.DoesBoxExist("box name")
	require.NoError(t, err)
	assert.False(t, exists)
	overlaps, err := tbdb.DoesSpanOverlap(0, 1)
	require.NoError(t, err)
	assert.False(t, overlaps)
	require.NoError(t, os.Remove(testdb))
}

func TestTBDB_AddBox(t *testing.T) {
	tbdb := setup(t)
	require.NoError(t, tbdb.AddBox("box-1", 1, 2))
	require.NoError(t, tbdb.AddBox("box-2", 1, 2))
	err := tbdb.AddBox("box-1", 3, 4)
	require.Error(t, err)
	assert.EqualError(t, err, "UNIQUE constraint failed: boxes.name")
	err = tbdb.AddBox("box-3", 5, 2)
	require.Error(t, err)
	assert.EqualError(t, err, "minTime is greater than maxTime")
}

func TestTBDB_AddSpan(t *testing.T) {
	tbdb := setup(t)
	require.NoError(t, tbdb.AddBox("box-1", 1, 2))
	require.NoError(t, tbdb.AddBox("box-2", 1, 2))
	require.NoError(t, tbdb.AddSpan(1, 10, "box-1"))
	require.NoError(t, tbdb.AddSpan(12, 14, "box-2"))
	err := tbdb.AddSpan(7, 11, "box-1")
	require.Error(t, err)
	assert.EqualError(t, err, "time overlaps existing span")
	err = tbdb.AddSpan(7, 11, "box-3")
	require.Error(t, err)
	assert.EqualError(t, err, "box box-3 doesn't exist")
	require.NoError(t, tbdb.AddSpan(15, 20, "box-1"))
	err = tbdb.AddSpan(2, 1, "box-2")
	require.Error(t, err)
	assert.EqualError(t, err, "start time is after end time")
	err = tbdb.AddSpan(time.Now().Unix()+86400, time.Now().Unix()+100000, "box-1")
	assert.EqualError(t, err, "time span is in the future")
}

func TestTBDB_GetBox(t *testing.T) {
	tbdb := setup(t)
	now := time.Now().Unix()
	tests := map[string]struct {
		name     string
		minTime  int64
		maxTime  int64
		create   bool
		expError bool
		errStr   string
	}{
		"good path":   {name: "box-1", minTime: 1, maxTime: 2, create: true, expError: false, errStr: ""},
		"invalid box": {name: "box-2", minTime: 1, maxTime: 2, create: false, expError: true, errStr: "sql: no rows in result set"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.create {
				require.NoError(t, tbdb.AddBox(tc.name, tc.minTime, tc.maxTime))
			}
			br, err := tbdb.GetBox(tc.name)
			if tc.expError {
				assert.EqualError(t, err, tc.errStr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.name, br.Name)
				assert.Equal(t, tc.minTime, br.MinTime)
				assert.Equal(t, tc.maxTime, br.MaxTime)
				assert.GreaterOrEqual(t, br.CreateTime, now)
			}
		})
	}
}

func TestTBDB_GetAllBoxes(t *testing.T) {
	now := time.Now().Unix()
	genBoxes := func(num int) []BoxRow {
		var boxes []BoxRow
		for i := 1; i <= num; i++ {
			boxes = append(boxes, BoxRow{Name: fmt.Sprintf("box-%d", i), MinTime: 5 * int64(i), MaxTime: 7 * int64(i)})
		}
		return boxes
	}
	tests := map[string]struct {
		boxes    []BoxRow
		expError bool
		errStr   string
	}{
		"good path": {boxes: genBoxes(3), expError: false, errStr: ""},
		"no boxes":  {boxes: genBoxes(0), expError: false, errStr: ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tbdb := setup(t)
			for _, box := range tc.boxes {
				require.NoError(t, tbdb.AddBox(box.Name, box.MinTime, box.MaxTime))
			}
			boxes, err := tbdb.GetAllBoxes()
			if tc.expError {
				assert.EqualError(t, err, tc.errStr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tc.boxes), len(boxes))
				for i := range boxes {
					assert.Equal(t, tc.boxes[i].Name, boxes[i].Name)
					assert.Equal(t, tc.boxes[i].MinTime, boxes[i].MinTime)
					assert.Equal(t, tc.boxes[i].MaxTime, boxes[i].MaxTime)
					assert.GreaterOrEqual(t, boxes[i].CreateTime, now)
				}
			}
		})
	}
}

func TestTBDB_GetSpansForBox(t *testing.T) {
	tbdb := setup(t)
	box := "box-1"
	ts := time.Now().Add(-365 * 24 * time.Hour).Unix()
	start := time.Now().Add(-1 * time.Hour)
	boxWithCreateTime(t, tbdb, box, 1, 2, ts)
	input := []SpanRow{
		{start.Unix(), start.Add(5 * time.Minute).Unix(), box},
		{start.Add(6 * time.Minute).Unix(), start.Add(7 * time.Minute).Unix(), box},
		{start.Add(9 * time.Minute).Unix(), start.Add(12 * time.Minute).Unix(), box},
	}
	for _, i := range input {
		require.NoError(t, tbdb.AddSpan(i.Start, i.End, i.Name))
	}
	spans, err := tbdb.GetSpansForBox("box-1")
	require.NoError(t, err)
	assert.Equal(t, len(input), len(spans))
	for i := range input {
		assert.Equal(t, input[i].Start, spans[i].Start)
		assert.Equal(t, input[i].End, spans[i].End)
		assert.Equal(t, input[i].Name, spans[i].Name)
	}
	require.NoError(t, tbdb.AddSpan(
		time.Now().Add(-368*24*time.Hour).Unix(), time.Now().Add(-367*24*time.Hour).Unix(), box))
	require.NoError(t, tbdb.AddSpan(
		time.Now().Add(-400*24*time.Hour).Unix(), time.Now().Add(-398*24*time.Hour).Unix(), box))
	spans, err = tbdb.GetSpansForBox("box-1")
	require.NoError(t, err)
	// shouldn't pick up last two spans which are older than the create time
	assert.Equal(t, len(input), len(spans))
}

func TestTBDB_UpdateBox(t *testing.T) {
	tbdb := setup(t)
	box := "box-1"
	require.NoError(t, tbdb.AddBox(box, 1, 2))
	require.NoError(t, tbdb.UpdateBox(box, 1, 3))
	br, err := tbdb.GetBox(box)
	require.NoError(t, err)
	assert.Equal(t, int64(3), br.MaxTime)
}

func TestTBDB_DeleteBox(t *testing.T) {
	tbdb := setup(t)
	box := "box-1"
	require.NoError(t, tbdb.AddBox(box, 1, 2))
	exists, err := tbdb.DoesBoxExist(box)
	require.NoError(t, err)
	assert.True(t, exists)
	require.NoError(t, tbdb.DeleteBox(box))
	exists, err = tbdb.DoesBoxExist(box)
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestTBDB_DeleteBoxAndSpans(t *testing.T) {
	tbdb := setup(t)
	box := "box-1"
	require.NoError(t, tbdb.AddBox(box, 1, 2))
	require.NoError(t, tbdb.AddSpan(1, 2, box))
	require.NoError(t, tbdb.AddSpan(5, 7, box))
	require.NoError(t, tbdb.AddBox("box-2", 1, 2))
	require.NoError(t, tbdb.AddSpan(8, 10, "box-2"))
	err := tbdb.DeleteBoxAndSpans(box)
	require.NoError(t, err)
	exists, err := tbdb.DoesBoxExist(box)
	require.NoError(t, err)
	assert.False(t, exists)
	_, err = tbdb.GetSpansForBox(box)
	assert.EqualError(t, err, "sql: no rows in result set")
}
