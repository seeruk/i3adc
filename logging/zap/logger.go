package zap

import (
	"github.com/seeruk/i3adc/logging"
	"go.uber.org/zap"
)

var _ logging.Logger = (*Logger)(nil)

// Logger implements this library's Logger interface using Uber's Zap logging library. Configuration
// of Zap is handled externally.
type Logger struct {
	SugaredLogger *zap.SugaredLogger
}

// NewLogger accepts a pre-configured Zap "sugared" logger, and returns a wrapped logger that
// satisfies this library's logging interface.
func NewLogger(sugaredLogger *zap.SugaredLogger) *Logger {
	return &Logger{
		SugaredLogger: sugaredLogger,
	}
}

// Debug logs a log entry at debug level.
func (l *Logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

// Debugf logs a log entry at debug level, using a given formatting string and arguments.
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.SugaredLogger.Debugf(msg, args...)
}

// Debugw logs a log entry at debug level, using the given message, and field pairs as arguments.
func (l *Logger) Debugw(msg string, args ...interface{}) {
	l.SugaredLogger.Debugw(msg, args...)
}

// Info logs a log entry at info level.
func (l *Logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

// Info logs a log entry at info level, using a given formatting string and arguments.
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.SugaredLogger.Infof(msg, args...)
}

// Infow logs a log entry at info level, using the given message, and field pairs as arguments.
func (l *Logger) Infow(msg string, args ...interface{}) {
	l.SugaredLogger.Infow(msg, args...)
}

// Warn logs a log entry at warning level.
func (l *Logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

// Warnf logs a log entry at warning level, using a given formatting string and arguments.
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.SugaredLogger.Warnf(msg, args...)
}

// Warnw logs a log entry at warning level, using the given message, and field pairs as arguments.
func (l *Logger) Warnw(msg string, args ...interface{}) {
	l.SugaredLogger.Warnw(msg, args...)
}

// Error logs a log entry at error level.
func (l *Logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

// Errorf logs a log entry at error level, using a given formatting string and arguments.
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.SugaredLogger.Errorf(msg, args...)
}

// Errorw logs a log entry at error level, using the given message, and field pairs as arguments.
func (l *Logger) Errorw(msg string, args ...interface{}) {
	l.SugaredLogger.Errorw(msg, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}

// Fatalf logs a log entry at fatal level, using a given formatting string and arguments.
func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.SugaredLogger.Fatalf(msg, args...)
}

// Fatalw logs a log fatal at error level, using the given message, and field pairs as arguments.
func (l *Logger) Fatalw(msg string, args ...interface{}) {
	l.SugaredLogger.Fatalw(msg, args...)
}

// With creates a logger based on this logger with the given field pairs added, so that any calls to
// the newly produced logger will include those fields.
func (l *Logger) With(args ...interface{}) logging.Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

// Sync attempts to flush the existing log entries synchronously using the underlying Zap logger.
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}
