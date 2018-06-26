package bolt

import (
	"github.com/gogo/protobuf/proto"
	"github.com/seeruk/i3adc/state"
)

// BucketOutputLayouts is the name of the storage bucket that output layouts are stored in.
const BucketOutputLayouts = "i3adc_output_layouts"

// TODO(seeruk): Do migrations in here, if they're necessary. For now, we only need to create a
// single bucket!

// Backend is a bolt-based implementation of i3adc's state backend interface.
type Backend struct {
	// TODO(seeruk): Needs a bolt DB instance?
}

// NewBackend returns a new bolt-based backend instance.
func NewBackend() *Backend {
	// TODO(seeruk): Here is probably where we'll try to create a bucket if it doesn't exist.
	return &Backend{}
}

// Read attempts to read a value from a bolt bucket under the given key into a given protobuf
// message.
func (b *Backend) Read(key string, val proto.Message) error {
	if key == "" {
		return state.ErrInvalidKey
	}

	if val == nil {
		return state.ErrInvalidValue
	}

	// TODO(seeruk): Get bytes.
	var value []byte

	return proto.Unmarshal(value, val)
}

// Write attempts to write the given protobuf message into a bolt bucket under the given key.
func (b *Backend) Write(key string, val proto.Message) error {
	if key == "" {
		return state.ErrInvalidKey
	}

	if val == nil {
		return state.ErrInvalidValue
	}

	value, err := proto.Marshal(val)
	if err != nil {
		return err
	}

	// TODO(seeruk): Use bytes.
	_ = value

	return nil
}

// Delete attempts to delete a key with the given name from the underlying bolt database.
func (b *Backend) Delete(key string) error {
	// TODO(seeruk): Do delete.

	return nil
}
