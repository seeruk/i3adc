package zap

import "github.com/seeruk/i3adc/internal/logging"

// Config contains all of the configuration relevant to a zap-based logger.
type Config struct {
	Level logging.Level `json:"level" consul:"level" env:"level"`
}
