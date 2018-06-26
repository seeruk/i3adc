package state

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// KeyLatestLayout is the key under which the most recently known layout hash is stored.
const KeyLatestLayout = "latest_layout"

var (
	// ErrInvalidKey is an error returned when a key is invalid.
	ErrInvalidKey = errors.New("state: invalid key")
	// ErrInvalidValue is an error returned when a value is invalid.
	ErrInvalidValue = errors.New("state: invalid value")
	// ErrNotFound is an error returned when a value is not found for a given key.
	ErrNotFound = errors.New("state: not found")
)

// Backend is an interface for interacting with application state at a basic level.
type Backend interface {
	Read(key string, val proto.Message) error
	Write(key string, val proto.Message) error
	Delete(key string) error
}
