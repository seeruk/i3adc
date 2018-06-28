package i3adc

import (
	"fmt"

	"github.com/seeruk/i3adc/logging"
	"github.com/seeruk/i3adc/logging/zap"
	"github.com/seeruk/i3adc/state/bolt"
	"github.com/seeruk/i3adc/xrandr"

	boltdb "github.com/coreos/bbolt"
)

// Resolver is a type that resolves application dependencies. Calling a resolver method is likely to
// call other resolver methods that are used to create the dependencies of the dependency you're
// asking for, and it will continue to follow that path until the entire dependency tree has been
// resolved for the dependency you're asking for.
type Resolver struct {
	boltDB       *boltdb.DB
	logger       logging.Logger
	xrandrClient *xrandr.Client
}

// NewResolver returns a new dependency resolver instance.
func NewResolver() *Resolver {
	resolver := &Resolver{}
	resolver.resolveEager()

	return resolver
}

// ResolverLogger resolves the singleton application logger instance.
func (r *Resolver) ResolveLogger() logging.Logger {
	if r.logger == nil {
		zapper := zap.New(zap.Config{
			Level: logging.DebugLevel,
		})

		r.logger = zap.NewLogger(zapper.Sugar())
	}

	return r.logger
}

// ResolveBoltDB resolves the singleton application bolt DB instance.
func (r *Resolver) ResolveBoltDB() *boltdb.DB {
	if r.boltDB == nil {
		db, err := bolt.OpenDB()
		if err != nil {
			panic(fmt.Sprintf("i3adc: failed to resolve bolt DB: %v", err))
		}

		r.boltDB = db
	}

	return r.boltDB
}

// ResolveStateBackend resolves a state backend instance, creating a new instance each time.
func (r *Resolver) ResolveStateBackend() *bolt.Backend {
	backend, err := bolt.NewBackend(r.ResolveBoltDB())
	if err != nil {
		panic(fmt.Sprintf("i3adc: failed to resolve state backend: %v", err))
	}

	return backend
}

// ResolveXrandrClient resolves the singleton application xrandr client instance.
func (r *Resolver) ResolveXrandrClient() *xrandr.Client {
	if r.xrandrClient == nil {
		client, err := xrandr.NewClient()
		if err != nil {
			panic(fmt.Sprintf("i3adc: failed to resolve xrandr client: %v", err))
		}

		r.xrandrClient = client
	}

	return r.xrandrClient
}

// resolveEager attempts to resolve dependencies that may error, so that those errors may be
// encountered at startup, instead of further into the application's life.
func (r *Resolver) resolveEager() {
	r.ResolveBoltDB()
	r.ResolveStateBackend()
	r.ResolveXrandrClient()
}
