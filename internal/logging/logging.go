package logging

// Log level constants. May be used for configuration later. Matches our logger interface's levels
// found below.
const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

// Level represents a logging level.
type Level string

// Logger is a general-purpose logging interface, based upon the interface provided by Uber's zap,
// allowing that library to be abstracted when needed.
type Logger interface {
	Debug(args ...interface{})
	Debugf(msg string, args ...interface{})
	Debugw(msg string, args ...interface{})

	Info(args ...interface{})
	Infof(msg string, args ...interface{})
	Infow(msg string, args ...interface{})

	Warn(args ...interface{})
	Warnf(msg string, args ...interface{})
	Warnw(msg string, args ...interface{})

	Error(args ...interface{})
	Errorf(msg string, args ...interface{})
	Errorw(msg string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(msg string, args ...interface{})
	Fatalw(msg string, args ...interface{})

	With(args ...interface{}) Logger

	// Sync tells the underlying log implementation to flush it's buffer. It may be useful to call
	// this before the application exits.
	Sync() error
}
