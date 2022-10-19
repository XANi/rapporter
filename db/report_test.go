package db

import (
	"github.com/efigence/go-mon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"os"
	"strconv"
	"testing"
)

func TestDB_AddReport(t *testing.T) {
	r := Report{
		Title:       "t-" + t.Name(),
		DeviceID:    "d-" + t.Name(),
		ComponentID: "c-" + t.Name(),
		State:       mon.StateOk,
		Content:     "**test report**",
	}
	db := DBTestInit(t)
	id, err := db.AddReport(r)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
	r.State = mon.StateCritical
	id, err = db.AddReport(r)
	require.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestDB_DeleteReport(t *testing.T) {
	db := DBTestInit(t)
	r := Report{
		Title:       "t-" + t.Name(),
		DeviceID:    "d-" + t.Name(),
		ComponentID: "c-" + t.Name(),
		Category:    "cat-" + t.Name(),
		State:       mon.StateOk,
		Content:     "**test report**",
		TTL:         600,
	}
	_, err := db.AddReport(r)
	require.NoError(t, err)
	reports, err := db.GetLatestReports()
	require.NoError(t, err)
	for id, _ := range reports {
		reports[id].CreatedAt = r.CreatedAt
		reports[id].UpdatedAt = r.UpdatedAt
	}
	assert.Contains(t, reports, r)
	err = db.DeleteReport(r.DeviceID, r.ComponentID)
	assert.NoError(t, err)
	reports, err = db.GetLatestReports()
	require.NoError(t, err)
	for id, _ := range reports {
		reports[id].CreatedAt = r.CreatedAt
		reports[id].UpdatedAt = r.UpdatedAt
	}
	assert.NotContains(t, reports, r)

}

func BenchmarkDB_AddReport(b *testing.B) {
	dsn := os.Getenv("TEST_DB_PATH")
	if len(dsn) < 1 {
		dsn = b.TempDir() + "/t.sqlite"
	}
	db, err := New(Config{
		DSN:    dsn,
		DbType: "sqlite",
		//DbType: "pgsql",
		Logger:   zaptest.NewLogger(b).Sugar(),
		testMode: true,
	})
	require.NoError(b, err)

	r := Report{
		Title:       "t-" + b.Name(),
		DeviceID:    "d-" + b.Name(),
		ComponentID: "c-" + b.Name(),
		Category:    "cat-" + b.Name(),
		State:       mon.StateOk,
		Content:     "**test report**",
		TTL:         600,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.DeviceID = r.DeviceID + strconv.Itoa(i)
		db.AddReport(r)
	}
}

func BenchmarkDB_UpdateReport(b *testing.B) {
	dsn := os.Getenv("TEST_DB_PATH")
	if len(dsn) < 1 {
		dsn = b.TempDir() + "/t.sqlite"
	}
	db, err := New(Config{
		DSN:    dsn,
		DbType: "sqlite",
		//DbType: "pgsql",
		Logger:   zaptest.NewLogger(b).Sugar(),
		testMode: true,
	})
	require.NoError(b, err)

	r := Report{
		Title:       "t-" + b.Name(),
		DeviceID:    "d-" + b.Name(),
		ComponentID: "c-" + b.Name(),
		Category:    "cat-" + b.Name(),
		State:       mon.StateOk,
		Content:     "**test report**",
		TTL:         600,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Content = "test-" + strconv.Itoa(i)
		db.AddReport(r)
	}
}
