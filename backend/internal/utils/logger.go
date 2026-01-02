package utils

import (
	"io"
	"os"
	"path/filepath"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger configures and returns a logrus logger based on the configuration
func SetupLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
		logger.Warnf("Invalid log level '%s', defaulting to 'info'", cfg.Log.Level)
	} else {
		logger.SetLevel(level)
	}

	// Set log format
	switch cfg.Log.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
		logger.Warnf("Invalid log format '%s', defaulting to 'text'", cfg.Log.Format)
	}

	// Set output destination
	switch cfg.Log.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "file":
		if cfg.Log.FilePath != "" {
			fileWriter := setupFileWriter(cfg)
			logger.SetOutput(fileWriter)
		} else {
			logger.SetOutput(os.Stdout)
			logger.Warn("Log file path not specified, falling back to stdout")
		}
	case "both":
		if cfg.Log.FilePath != "" {
			fileWriter := setupFileWriter(cfg)
			multiWriter := io.MultiWriter(os.Stdout, fileWriter)
			logger.SetOutput(multiWriter)
		} else {
			logger.SetOutput(os.Stdout)
			logger.Warn("Log file path not specified, falling back to stdout only")
		}
	default:
		logger.SetOutput(os.Stdout)
		logger.Warnf("Invalid log output '%s', defaulting to stdout", cfg.Log.Output)
	}

	return logger
}

// setupFileWriter creates a file writer with log rotation using lumberjack
func setupFileWriter(cfg *config.Config) io.Writer {
	// Ensure the log directory exists
	logDir := filepath.Dir(cfg.Log.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Errorf("Failed to create log directory '%s': %v", logDir, err)
		return os.Stdout
	}

	return &lumberjack.Logger{
		Filename:   cfg.Log.FilePath,
		MaxSize:    cfg.Log.MaxSize,    // megabytes
		MaxBackups: cfg.Log.MaxBackups, // number of backups
		MaxAge:     cfg.Log.MaxAge,     // days
		Compress:   cfg.Log.Compress,   // compress old files
		LocalTime:  true,               // use local time for backup file names
	}
}

// GetLoggerWithFields creates a logger with predefined fields
func GetLoggerWithFields(logger *logrus.Logger, fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// GetRequestLogger creates a logger with request-specific fields
func GetRequestLogger(logger *logrus.Logger, requestID, method, path string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"method":     method,
		"path":       path,
	})
}
