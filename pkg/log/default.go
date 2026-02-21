package ravandlog

import (
	"context"

	"github.com/sirupsen/logrus"
)

var defaultLogger *Logger

func init() {
	defaultLogger, _ = NewLogger(Config{
		Level: "debug",
	})
}

func SetupDefaultLogger(l *Logger) {
	defaultLogger = l
}

func GetLogger() *Logger {
	return defaultLogger
}

func Trace(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Trace(args...)
}

func Debug(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Debug(args...)
}

func Info(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Info(args...)
}

func Warn(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Warn(args...)
}

func Error(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Error(args...)
}

func Fatal(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Fatal(args...)
}

func Debugln(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Debugln(args...)
}

func Infoln(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Infoln(args...)
}

func Warnln(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Warnln(args...)
}

func Errorln(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	ctx := GetCurrentFunctionName(context.Background())
	defaultLogger.WithContext(ctx).WithFields(defaultLogger.extractFieldsFromContext(ctx)).entry.Fatalln(args...)
}

func With(key string, value interface{}) *logrus.Entry {
	return defaultLogger.l.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return defaultLogger.l.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return defaultLogger.l.WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return defaultLogger.l.WithContext(ctx)
}
