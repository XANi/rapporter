package db

import (
	"fmt"
	"github.com/efigence/go-mon"
	"time"
)

type Report struct {
	Title       string `gorm:"length:2048" json:"title"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Category    string    `json:"category" gorm:"length:255"`
	DeviceID    string    `gorm:"primaryKey;length:255" json:"device_id"`
	ComponentID string    `gorm:"primaryKey;length:255" json:"component_id"`
	TTL         uint      `json:"ttl"`
	State       mon.State `json:"state"`
	Content     string    `json:"content"`
}

func (r *Report) Validate() error {
	if len(r.Title) == 0 || len(r.Title) > 2048 {
		return fmt.Errorf("title needs to be between 1 and 2048 characters")
	}
	if len(r.DeviceID) == 0 || len(r.DeviceID) > 2048 {
		return fmt.Errorf("device_id needs to be between 1 and 2048 characters")
	}
	if len(r.ComponentID) == 0 || len(r.ComponentID) > 2048 {
		return fmt.Errorf("component_id needs to be between 1 and 2048 characters")
	}
	if len(r.Category) > 255 {
		return fmt.Errorf("categoru needs to be shorter than 255 characters")
	}
	if r.State < 1 || r.State > mon.StateUnknown {
		return fmt.Errorf("state must be between 1(ok) and 4(unknown)")
	}
	return nil
}

type ReportState struct {
	Report    Report `gorm:"embedded"`
	ExpiresAt time.Time
}

type ReportHistory struct {
	ID     uint   `gorm:"primarykey"`
	Report Report `gorm:"embedded"`
}
type ReportHistoryMigrate struct {
	ID          uint   `gorm:"primarykey"`
	Report      Report `gorm:"embedded"`
	DeviceID    string `gorm:"length:255"`
	ComponentId string `gorm:"length:255"`
}

func (ReportHistoryMigrate) TableName() string {
	return "report_history"
}
