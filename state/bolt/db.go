package bolt

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/bbolt"
	"github.com/seeruk/i3adc/state"
)

// OpenDB attempts to open the bolt-based database, returning a new bolt DB instance.
func OpenDB() (*bolt.DB, error) {
	localDir, err := state.LocalDirectory()
	if err != nil {
		return nil, fmt.Errorf("bolt: failed to get local directory: %v", err)
	}

	// Make the local directory if it does not exist.
	err = os.MkdirAll(localDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Open / create as 0700 so that we (and only we) have read/write access.
	return bolt.Open(filepath.Join(localDir, "i3adc.db"), 0600, &bolt.Options{
		Timeout: 5 * time.Second,
	})
}
