package db

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

type dbLogger struct {
	l *zap.SugaredLogger
}

func (t *dbLogger) LogMode(level logger.LogLevel) logger.Interface {
	return t
}
func (t *dbLogger) Info(ctx context.Context, fmt string, d ...interface{}) {
	t.l.Infof(fmt, d...)

}
func (t *dbLogger) Warn(ctx context.Context, fmt string, d ...interface{}) {
	t.l.Warnf(fmt, d...)

}
func (t *dbLogger) Error(ctx context.Context, fmt string, d ...interface{}) {
	t.l.Errorf(fmt, d...)
}
func (t *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

}
