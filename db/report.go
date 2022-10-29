package db

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func (db *DB) AddReport(r Report) (id int, err error) {
	if r.TTL == 0 {
		r.TTL = uint(db.cfg.DefaultExpiryInterval.Seconds())
	}
	if db.cfg.MaxExpiryInterval > 0 && r.TTL > uint(db.cfg.MaxExpiryInterval.Seconds()) {
		r.TTL = uint(db.cfg.MaxExpiryInterval.Seconds())
	}
	expire := uint(db.cfg.MaxExpiryInterval.Seconds())
	invalid := r.TTL
	if r.ExpireIn > 0 {
		expire = r.ExpireIn
	}
	reportS := ReportState{
		ExpiresAt: time.Now().Add(time.Duration(expire) * time.Second),
		InvalidAt: time.Now().Add(time.Duration(invalid) * time.Second),
		Report:    r,
	}
	reportH := ReportHistory{
		Report: r,
	}
	tx := db.d.Begin()
	defer tx.Rollback()
	if q := tx.Save(&reportS); q.Error != nil {
		return 0, q.Error
	}
	q := tx.Create(&reportH)
	if q.Error != nil {
		return 0, q.Error
	}
	tx.Commit()
	return int(reportH.ID), tx.Error
}

func (db *DB) DeleteReport(deviceID string, componentID string) error {
	q := db.d.Delete(&ReportState{}, &Report{DeviceID: deviceID, ComponentID: componentID})

	if q.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return q.Error
}

func (db *DB) GetLatestReports() ([]Report, error) {
	r := []Report{}
	q := db.d.Model(&ReportState{}).Find(&r)
	if q.Error == gorm.ErrRecordNotFound {
		return r, nil
	}
	return r, q.Error
}
func (db *DB) GetReport(deviceId string, componentId string) (Report, error) {
	r := Report{}
	q := db.d.Model(&ReportState{}).Limit(1).Find(&r, &Report{
		DeviceID:    deviceId,
		ComponentID: componentId,
	})
	if q.Error == gorm.ErrRecordNotFound || q.RowsAffected == 0 {
		return r, fmt.Errorf("not found")
	}
	return r, q.Error
}
