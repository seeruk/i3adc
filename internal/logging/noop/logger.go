package noop

import "github.com/seeruk/i3adc/internal/logging"

// Logger is a no-op logger that implements the interface of this package, but doesn't actually
// do anything... at all.
type Logger struct{}

// NewLogger creates a new no-op logger instance.
func NewLogger() *Logger {
	return &Logger{}
}

// Debug is a no-op function.
func (l *Logger) Debug(_ ...interface{}) {}

// Debugf is a no-op function.
func (l *Logger) Debugf(_ string, _ ...interface{}) {}

// Debugw is a no-op function.
func (l *Logger) Debugw(_ string, _ ...interface{}) {}

// Info is a no-op function.
func (l *Logger) Info(_ ...interface{}) {}

// Infof is a no-op function.
func (l *Logger) Infof(_ string, _ ...interface{}) {}

// Infow is a no-op function.
func (l *Logger) Infow(_ string, _ ...interface{}) {}

// Warn is a no-op function.
func (l *Logger) Warn(_ ...interface{}) {}

// Warnf is a no-op function.
func (l *Logger) Warnf(_ string, _ ...interface{}) {}

// Warnw is a no-op function.
func (l *Logger) Warnw(_ string, _ ...interface{}) {}

// Error is a no-op function.
func (l *Logger) Error(_ ...interface{}) {}

// Errorf is a no-op function.
func (l *Logger) Errorf(_ string, _ ...interface{}) {}

// Errorw is a no-op function.
func (l *Logger) Errorw(_ string, _ ...interface{}) {}

// Fatal is a no-op function.
func (l *Logger) Fatal(_ ...interface{}) {}

// Fatalf is a no-op function.
func (l *Logger) Fatalf(_ string, _ ...interface{}) {}

// Fatalw is a no-op function.
func (l *Logger) Fatalw(_ string, _ ...interface{}) {}

// With is a no-op function.
func (l *Logger) With(_ ...interface{}) logging.Logger {
	return l
}

// Sync is a no-op function.
func (l *Logger) Sync() error {
	return nil
}
