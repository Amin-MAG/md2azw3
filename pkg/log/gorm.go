package ravandlog

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	*Logger
}

// LogMode sets the log mode for the GormLogger
func (gl GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	var logLevel string
	switch level {
	case logger.Silent:
		logLevel = "debug"
	case logger.Info:
		logLevel = "info"
	case logger.Warn:
		logLevel = "warn"
	case logger.Error:
		logLevel = "error"
	}

	l, err := logrus.ParseLevel(logLevel)
	if err == nil {
		gl.l.SetLevel(l)
	}

	return gl
}

// Debug logs an debugging message
func (gl GormLogger) Debug(ctx context.Context, s string, i ...interface{}) {
	gl.WithContext(ctx).Debugf(ctx, s, i...)
}

// Info logs an informational message
func (gl GormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	gl.WithContext(ctx).Infof(ctx, s, i...)
}

// Warn logs a warning message
func (gl GormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	gl.WithContext(ctx).Warnf(ctx, s, i...)
}

// Error logs an error message
func (gl GormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	gl.WithContext(ctx).Errorf(ctx, s, i...)
}

// Trace logs a SQL statement with its duration
func (gl GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, rows := fc()
		gl.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"sql":        sql,
			"rows":       rows,
			"elapsed_ms": time.Since(begin).Milliseconds(),
		}).Error(ctx, "SQL execution failed")
	} else {
		sql, rows := fc()
		gl.WithContext(ctx).WithFields(logrus.Fields{
			"sql":        sql,
			"rows":       rows,
			"elapsed_ms": time.Since(begin).Milliseconds(),
		}).Info(ctx, "SQL executed")
	}
}
