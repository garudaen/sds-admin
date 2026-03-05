package logger

import (
	"io"
	"os"
	"path/filepath"

	"sds-admin/internal/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *logrus.Logger

// Init initializes the global logger with configuration.
// Init 使用配置初始化全局日志记录器
func Init(cfg *config.LogConfig) error {
	log = logrus.New()

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	switch cfg.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			DisableColors:   true,
		})
	}

	var writer io.Writer
	switch cfg.Output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		if err := os.MkdirAll(filepath.Dir(cfg.Output), 0755); err != nil {
			return err
		}
		writer = &lumberjack.Logger{
			Filename:   cfg.Output,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
	}

	log.SetOutput(writer)

	return nil
}

// GetLogger returns the global logger instance.
// GetLogger 返回全局日志记录器实例
func GetLogger() *logrus.Logger {
	if log == nil {
		log = logrus.New()
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			DisableColors:   true,
		})
	}
	return log
}

// Debug logs a message at level Debug.
// Debug 在Debug级别记录消息
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf logs a formatted message at level Debug.
// Debugf 在Debug级别记录格式化消息
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info logs a message at level Info.
// Info 在Info级别记录消息
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof logs a formatted message at level Info.
// Infof 在Info级别记录格式化消息
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn logs a message at level Warn.
// Warn 在Warn级别记录消息
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf logs a formatted message at level Warn.
// Warnf 在Warn级别记录格式化消息
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error logs a message at level Error.
// Error 在Error级别记录消息
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf logs a formatted message at level Error.
// Errorf 在Error级别记录格式化消息
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1.
// Fatal 在Fatal级别记录消息，然后进程将以状态1退出
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf logs a formatted message at level Fatal then the process will exit with status set to 1.
// Fatalf 在Fatal级别记录格式化消息，然后进程将以状态1退出
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithField creates an entry with a single field.
// WithField 创建带有单个字段的日志条目
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields creates an entry with multiple fields.
// WithFields 创建带有多个字段的日志条目
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}
