package zap

import (
	"os"
	"syscall"

	"github.com/seeruk/i3adc/internal/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

// levelMap maps icelog's levels to uber-go/zap's levels.
var levelMap = map[logging.Level]zapcore.Level{
	logging.DebugLevel: zapcore.DebugLevel,
	logging.InfoLevel:  zapcore.InfoLevel,
	logging.WarnLevel:  zapcore.WarnLevel,
	logging.ErrorLevel: zapcore.ErrorLevel,
	logging.FatalLevel: zapcore.FatalLevel,
}

// New returns a new uber-go/zap instance, with some sensible automation around it's configuration
// set up.
func New(config Config) *zap.Logger {
	// If we are running in production, stdin will not be a terminal. Otherwise, we should use a
	// more friendly looking output style.
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	if isTerminal() {
		encoderConf := zap.NewDevelopmentEncoderConfig()
		encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConf.MessageKey = "message"

		encoder = zapcore.NewConsoleEncoder(encoderConf)
	}

	writer := os.Stdout

	minLevel, ok := levelMap[config.Level]
	if !ok {
		minLevel = zapcore.InfoLevel
	}

	enabler := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= minLevel
	})

	return zap.New(zapcore.NewCore(encoder, writer, enabler))
}

// isTerminal will return true if the application's stdin appears to be a terminal.
func isTerminal() bool {
	return terminal.IsTerminal(syscall.Stdin)
}
