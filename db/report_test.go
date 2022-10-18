package db

import (
	"github.com/efigence/go-mon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	}
	assert.Contains(t, reports, r)
	err = db.DeleteReport(r.DeviceID, r.ComponentID)
	assert.NoError(t, err)
	reports, err = db.GetLatestReports()
	require.NoError(t, err)
	for id, _ := range reports {
		reports[id].CreatedAt = r.CreatedAt
	}
	assert.NotContains(t, reports, r)

}
