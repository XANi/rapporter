package db

import (
	"github.com/efigence/go-mon"
	"time"
)

type Report struct {
	Title       string `gorm:"length:2048"`
	DeviceID    string `gorm:"primaryKey;length:255"`
	ComponentId string `gorm:"primaryKey;length:255"`
	TTL         uint
	State       mon.State
	Content     string
}

type ReportState struct {
	Report    Report `gorm:"embedded"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type ReportHistory struct {
	ID        uint   `gorm:"primarykey"`
	Report    Report `gorm:"embedded"`
	CreatedAt time.Time
}
type ReportHistoryMigrate struct {
	ID          uint   `gorm:"primarykey"`
	Report      Report `gorm:"embedded"`
	CreatedAt   time.Time
	DeviceID    string `gorm:"length:255"`
	ComponentId string `gorm:"length:255"`
}

func (ReportHistoryMigrate) TableName() string {
	return "report_history"
}
