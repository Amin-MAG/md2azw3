package ravandlog

import (
	"context"
	"fmt"
	"io"
	. "os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	DefaultMaxNumberOfFiles = 10
	DefaultMaxSizeInBytes   = 100000
)

// Config represents details about a Logger
type Config struct {
	Level            string
	IsPrettyPrint    bool
	IsReportCaller   bool
	OutputFileConfig *OutputFileConfig
}

// Entry is a wrapper for logrus.Entry
type Entry struct {
	logger *Logger
	entry  *logrus.Entry
}

func (l *Entry) Trace(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Trace(args...)
}

func (l *Entry) Debug(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Debug(args...)
}

func (l *Entry) Info(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Info(args...)
}

func (l *Entry) Warn(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Warn(args...)
}

func (l *Entry) Error(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Error(args...)
}

func (l *Entry) Fatal(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Fatal(args...)
}

func (l *Entry) Debugf(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Debugf(s, args...)
}

func (l *Entry) Infof(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Infof(s, args...)
}

func (l *Entry) Warnf(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Warnf(s, args...)
}

func (l *Entry) Errorf(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.entry.WithContext(ctx).WithFields(l.logger.extractFieldsFromContext(ctx)).Errorf(s, args...)
}

func (l *Entry) With(key string, value interface{}) *Entry {
	return &Entry{
		entry:  l.entry.WithField(key, value),
		logger: l.logger,
	}
}

func (l *Entry) WithFields(fields logrus.Fields) *Entry {
	return &Entry{
		entry:  l.entry.WithFields(fields),
		logger: l.logger,
	}
}

func (l *Entry) WithError(err error) *Entry {
	return &Entry{
		entry:  l.entry.WithError(err),
		logger: l.logger,
	}
}

func (l *Entry) WithTime(t time.Time) *Entry {
	return &Entry{
		entry:  l.entry.WithTime(t),
		logger: l.logger,
	}
}

func (l *Entry) WithContext(ctx context.Context) *Entry {
	ctx = GetCurrentFunctionName(ctx)
	return &Entry{
		entry:  l.entry.WithContext(ctx),
		logger: l.logger,
	}
}

// Logger is a wrapper for logrus.Logger.
type Logger struct {
	l      *logrus.Logger
	config Config
}

// NewLogger creates new Logger with Config.
func NewLogger(config Config) (*Logger, error) {
	// Create new logrus log
	l := logrus.New()
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("the name for log level is not valid: %s", err)
	}
	l.SetLevel(level)
	if config.IsPrettyPrint {
		l.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			PrettyPrint:     true,
		})
	} else {
		l.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}
	// Check if it is going to be saved on file
	if config.OutputFileConfig != nil {
		// Default Values
		if config.OutputFileConfig.MaxNumberOfFiles == 0 {
			config.OutputFileConfig.MaxNumberOfFiles = DefaultMaxNumberOfFiles
		}
		if config.OutputFileConfig.MaxSizeInBytes == 0 {
			config.OutputFileConfig.MaxSizeInBytes = DefaultMaxSizeInBytes
		}

		// Set the file as an output
		if f, err := OpenFile((*config.OutputFileConfig).FullPath(0), O_RDWR|O_CREATE|O_APPEND, 0666); err == nil {
			l.SetOutput(io.MultiWriter(Stdout, f))
		} else {
			return nil, fmt.Errorf("cannot open log file: %s", err)
		}

		// Hook for checking storage limitation.
		l.AddHook(NewFileThresholdHook(*config.OutputFileConfig))
	}

	return &Logger{
		l:      l,
		config: config,
	}, nil
}

func (l *Logger) CloneGormLogger() (*GormLogger, error) {
	// Create a new similar Logger based on the configuration
	newLogger, err := NewLogger(
		Config{
			Level:          l.config.Level,
			IsPrettyPrint:  l.config.IsPrettyPrint,
			IsReportCaller: l.config.IsReportCaller,
		})
	if err != nil {
		return nil, err
	}

	return &GormLogger{
		newLogger,
	}, nil
}

func (l *Logger) Trace(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Trace(args...)
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Debug(args...)
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Info(args...)
}

func (l *Logger) Warn(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Warn(args...)
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Error(args...)
}

func (l *Logger) Fatal(ctx context.Context, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Fatal(args...)
}

func (l *Logger) Debugf(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Debugf(s, args...)
}

func (l *Logger) Infof(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Infof(s, args...)
}

func (l *Logger) Warnf(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Warnf(s, args...)
}

func (l *Logger) Errorf(ctx context.Context, s string, args ...interface{}) {
	ctx = GetCurrentFunctionName(ctx)
	l.l.WithContext(ctx).WithFields(l.extractFieldsFromContext(ctx)).Errorf(s, args...)
}

func (l *Logger) With(key string, value interface{}) *Entry {
	return &Entry{
		entry:  l.l.WithField(key, value),
		logger: l,
	}
}

func (l *Logger) WithError(err error) *Entry {
	return &Entry{
		entry:  l.l.WithError(err),
		logger: l,
	}
}

func (l *Logger) WithTime(t time.Time) *Entry {
	return &Entry{
		entry:  l.l.WithTime(t),
		logger: l,
	}
}

func (l *Logger) WithContext(ctx context.Context) *Entry {
	return &Entry{
		entry:  l.l.WithContext(ctx),
		logger: l,
	}
}
