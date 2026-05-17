// Package logger provides a centralized logrus initialization for the application.
// Log level is controlled via the LOG_LEVEL environment variable (default: "error").
// Supported levels: debug, info, warn, error.
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Level is a custom log-level type backed by an iota so callers get
// compile-time safety instead of raw strings.
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

// ParseLevel converts a LOG_LEVEL string (e.g. from an env var) into the
// custom Level type. Unrecognized values default to Error and emit a warning.
func ParseLevel(s string) Level {
	switch s {
	case "debug", "DEBUG":
		return Debug
	case "info", "INFO":
		return Info
	case "warn", "warning", "WARN", "WARNING":
		return Warn
	case "error", "ERROR":
		return Error
	default:
		logrus.Warnf("Unknown LOG_LEVEL %q, defaulting to 'error'", s)
		return Error
	}
}

// toLogrusLevel maps the custom Level to the equivalent logrus.Level.
func toLogrusLevel(l Level) logrus.Level {
	switch l {
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warn:
		return logrus.WarnLevel
	default: // Error and any future unknown values
		return logrus.ErrorLevel
	}
}

// Init configures the global logrus logger with a human-readable TextFormatter.
// Call it twice in main: once with logger.Error before config loads, then again
// with cfg.LogLevel after so the user-configured level takes effect.
func Init(level Level) {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)

	resolved := toLogrusLevel(level)
	logrus.SetLevel(resolved)
	logrus.WithField("level", resolved.String()).Info("Logger initialized")
}
