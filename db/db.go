package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Config struct {
	DSN                   string
	DbType                string
	Logger                *zap.SugaredLogger
	DefaultExpiryInterval time.Duration `yaml:"default_expiry"`
	MaxExpiryInterval     time.Duration `yaml:"max_expiry"`
	HistoryDays           int           `yaml:"history_days"`
	testMode              bool
}

type DB struct {
	d   *gorm.DB
	l   *zap.SugaredLogger
	cfg Config
}

func New(cfg Config) (*DB, error) {
	var dbConn *gorm.DB
	var err error
	dbCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger: &dbLogger{cfg.Logger},
	}
	switch cfg.DbType {
	case "sqlite":
		dbConn, err = gorm.Open(sqlite.Open(cfg.DSN+"?_journal_mode=WAL&_synchronous=NORMAL"), dbCfg)
	case "pgsql":
		dbConn, err = gorm.Open(postgres.Open(cfg.DSN), dbCfg)
	default:
		return nil, fmt.Errorf("db type [%s] not supported", err)
	}
	if err != nil {
		return &DB{}, err
	}
	if cfg.DefaultExpiryInterval == 0 {
		cfg.DefaultExpiryInterval = time.Hour * 24 * 8
	}
	if cfg.MaxExpiryInterval == 0 {
		cfg.MaxExpiryInterval = time.Hour * 24 * 30
	}
	if cfg.HistoryDays == 0 {
		cfg.HistoryDays = 365 * 5
	}
	var dbObj DB
	if cfg.Logger != nil {
		dbObj.l = cfg.Logger
	} else {
		dbObj.l = zap.S()
	}
	dbObj.d = dbConn
	dbObj.cfg = cfg
	migrations := append(make([]interface{}, 0),
		ReportState{},
		ReportHistoryMigrate{},
	)
	for _, table := range migrations {
		err = dbObj.d.AutoMigrate(table)
		if err != nil {
			return nil, fmt.Errorf("error migrating %T: %s", table, err)
		}
	}
	if !cfg.testMode {
		go dbObj.cleaner()
	}
	return &dbObj, err
}

func (db *DB) cleaner() {
	db.l.Infof("cleaner set to remove reports older than %d days", db.cfg.HistoryDays)
	for {
		q := db.d.Exec("DELETE FROM report_history WHERE created_at < ?",
			time.Now().Add(-1*time.Hour*time.Duration(24*db.cfg.HistoryDays)))
		if q.RowsAffected > 0 {
			db.l.Infof("deleted %d old records from history", q.RowsAffected)
		}
		if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
			db.l.Warnf("error running cleanup: %s", q.Error)
		}
		q = db.d.Exec("DELETE FROM report_state WHERE expires_at < ?", time.Now())
		if q.RowsAffected > 0 {
			db.l.Infof("deleted %d expired records", q.RowsAffected)
		}
		if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
			db.l.Warnf("error running cleanup: %s", q.Error)
		}
		q = db.d.Exec("UPDATE report_state SET state = 4 WHERE invalid_at < ?", time.Now())
		if q.RowsAffected > 0 {
			db.l.Infof("deleted %d expired records", q.RowsAffected)
		}
		if q.Error != nil && q.Error != gorm.ErrRecordNotFound {
			db.l.Warnf("error running cleanup: %s", q.Error)
		}

		time.Sleep(time.Hour * 10)
	}
}
