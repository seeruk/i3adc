package state

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// LocalDirectory returns the directory that i3adc will use to store all configuration / data for
// the current user.
func LocalDirectory() (string, error) {
	home, err := homeDirectory()
	if err != nil {
		return "", fmt.Errorf("state: couldn't get local directory to use: %v", err)
	}

	return filepath.Join(home, ".i3adc"), nil
}

// homeDirectory attempts to get the current user's home directory.
func homeDirectory() (string, error) {
	// Try to use the environment first.
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	if usr.HomeDir == "" {
		return "", errors.New("invalid (empty) home directory path")
	}

	return usr.HomeDir, nil
}
