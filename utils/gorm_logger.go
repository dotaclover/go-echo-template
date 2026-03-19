package utils

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogrusLogger 使用 logrus 实现 GORM 日志接口
type GormLogrusLogger struct {
	config logger.Config
}

func NewGormLogrusLogger(cfg logger.Config) logger.Interface {
	return &GormLogrusLogger{config: cfg}
}

func (l *GormLogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	copy := *l
	copy.config.LogLevel = level
	return &copy
}

func (l *GormLogrusLogger) shouldLog(level logger.LogLevel) bool {
	return level <= l.config.LogLevel
}

func (l *GormLogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if !l.shouldLog(logger.Info) {
		return
	}
	Logger.WithField("component", "gorm").Infof(msg, data...)
}

func (l *GormLogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if !l.shouldLog(logger.Warn) {
		return
	}
	Logger.WithField("component", "gorm").Warnf(msg, data...)
}

func (l *GormLogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if !l.shouldLog(logger.Error) {
		return
	}
	Logger.WithField("component", "gorm").Errorf(msg, data...)
}

func (l *GormLogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if !l.shouldLog(logger.Info) {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := logrus.Fields{
		"component":     "gorm",
		"elapsed_ms":    elapsed.Milliseconds(),
		"rows_affected": rows,
		"sql":           sql,
	}

	if err != nil {
		if l.config.IgnoreRecordNotFoundError && errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
		Logger.WithFields(fields).WithError(err).Error("sql_error")
		return
	}

	if l.config.SlowThreshold > 0 && elapsed > l.config.SlowThreshold {
		Logger.WithFields(fields).Warn("slow_sql")
		return
	}

	if l.config.LogLevel >= logger.Info {
		Logger.WithFields(fields).Debug("sql")
	}
}
