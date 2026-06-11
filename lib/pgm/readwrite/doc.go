// Package readwrite provides readers and writers for probabilistic model
// file formats: BIF, XMLBIF, NET, UAI, XDSL, PomdpX, XBN, CSV, JSON,
// and datascience-native XML.
package readwrite

import (
	"fmt"
	"io"
)

// MaxInputSize is the default maximum number of bytes accepted by any reader
// in this package. Inputs exceeding this limit are rejected to prevent
// denial-of-service via oversized or maliciously crafted files (e.g., XML bombs).
// Use the *WithLimit variants (e.g., ReadXMLBIFWithLimit) for larger models.
const MaxInputSize = 1 << 20 // 1 MB

// readLimited reads up to MaxInputSize bytes from r.
func readLimited(r io.Reader) ([]byte, error) {
	return readLimitedN(r, MaxInputSize)
}

// readLimitedN reads up to maxBytes from r. If the input exceeds maxBytes,
// it returns an error without buffering the entire stream.
func readLimitedN(r io.Reader, maxBytes int) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, fmt.Errorf("readwrite: maxBytes must be positive, got %d", maxBytes)
	}
	limit := int64(maxBytes) + 1
	data, err := io.ReadAll(io.LimitReader(r, limit))
	if err != nil {
		return nil, err
	}
	if len(data) > maxBytes {
		return nil, fmt.Errorf("readwrite: input exceeds maximum size (%d bytes)", maxBytes)
	}
	return data, nil
}
