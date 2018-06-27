package bolt

import (
	"github.com/coreos/bbolt"
	"github.com/seeruk/i3adc/state"
)

// BucketOutputLayouts is the name of the storage bucket that output layouts are stored in.
const BucketOutputLayouts = "i3adc_output_layouts"

// Backend is a bolt-based implementation of i3adc's state backend interface.
type Backend struct {
	db *bolt.DB
}

// NewBackend returns a new bolt-based backend instance.
func NewBackend(db *bolt.DB) (*Backend, error) {
	// Always ensure that the bucket exists.
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketOutputLayouts))
		return err
	})

	if err != nil {
		return nil, err
	}

	return &Backend{
		db: db,
	}, nil
}

// Read attempts to read a value from a bolt bucket under the given key, returning the raw bytes.
func (b *Backend) Read(key string) ([]byte, error) {
	var bs []byte
	if key == "" {
		return bs, state.ErrInvalidKey
	}

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketOutputLayouts))
		bs = bucket.Get([]byte(key))

		return nil
	})

	return bs, err
}

// Write attempts to write the given bytes into a bolt bucket under the given key.
func (b *Backend) Write(key string, val []byte) error {
	if key == "" {
		return state.ErrInvalidKey
	}

	if val == nil {
		return state.ErrInvalidValue
	}

	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketOutputLayouts))
		return bucket.Put([]byte(key), val)
	})
}

// Delete attempts to delete a key with the given name from the underlying bolt database.
func (b *Backend) Delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketOutputLayouts))
		return bucket.Delete([]byte(key))
	})
}
