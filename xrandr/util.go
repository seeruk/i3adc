package xrandr

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// calculateHashForOutputs takes a set of outputs and produces an MD5 hash of all of the properties
// of the connected outputs. This serves as a way of uniquely identifying a set of outputs.
func calculateHashForOutputs(outputs []Output) (string, error) {
	sum := md5.New()

	for _, output := range outputs {
		if !output.IsConnected {
			continue
		}

		io.WriteString(sum, output.Name)
		sum.Write(output.Properties["EDID"])
	}

	return hex.EncodeToString(sum.Sum(nil)), nil
}
