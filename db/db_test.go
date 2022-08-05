package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const dbName = "test.db"

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
	tempDir, err := ioutil.TempDir(os.TempDir(), "timebox")
	require.NoError(t, err)
	testdb := filepath.Join(tempDir, filepath.FromSlash(dbName))
	os.Remove(testdb)
	tbdb := NewDBWithName(testdb)
	require.NoError(t, tbdb.CreateDB())
	require.NoError(t, tbdb.AddBox("box-1", 1, 2))
	require.NoError(t, tbdb.AddBox("box-2", 1, 2))
	err = tbdb.AddBox("box-1", 3, 4)
	require.Error(t, err)
	assert.EqualError(t, err, "UNIQUE constraint failed: boxes.name")
	err = tbdb.AddBox("box-3", 5, 2)
	require.Error(t, err)
	assert.EqualError(t, err, "minTime is greater than maxTime")
}

func TestTBDB_AddSpan(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "timebox")
	require.NoError(t, err)
	testdb := filepath.Join(tempDir, filepath.FromSlash(dbName))
	os.Remove(testdb)
	tbdb := NewDBWithName(testdb)
	require.NoError(t, tbdb.CreateDB())
	require.NoError(t, tbdb.AddBox("box-1", 1, 2))
	require.NoError(t, tbdb.AddBox("box-2", 1, 2))
	require.NoError(t, tbdb.AddSpan(1, 10, "box-1"))
	require.NoError(t, tbdb.AddSpan(12, 14, "box-2"))
	err = tbdb.AddSpan(7, 11, "box-1")
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
