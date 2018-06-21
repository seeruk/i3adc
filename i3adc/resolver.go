package internal

import (
	"github.com/seeruk/i3adc/logging"
	"github.com/seeruk/i3adc/logging/zap"
)

// Resolver is a type that resolves application dependencies. Calling a resolver method is likely to
// call other resolver methods that are used to create the dependencies of the dependency you're
// asking for, and it will continue to follow that path until the entire dependency tree has been
// resolved for the dependency you're asking for.
type Resolver struct {
	logger logging.Logger
}

// NewResolver returns a new dependency resolver instance.
func NewResolver() *Resolver {
	return &Resolver{}
}

// ResolverLogger resolves the singleton application logger instance.
func (r *Resolver) ResolveLogger() logging.Logger {
	zapper := zap.New(zap.Config{
		Level: logging.DebugLevel,
	})

	return zap.NewLogger(zapper.Sugar())
}
