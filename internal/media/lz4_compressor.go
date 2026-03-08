package media

import (
	"fmt"
	"github.com/pierrec/lz4/v4"
)

// LZ4Compressor implements the domain.Compressor interface using LZ4.
type LZ4Compressor struct{}

// NewLZ4Compressor creates a new LZ4Compressor.
func NewLZ4Compressor() *LZ4Compressor {
	return &LZ4Compressor{}
}

// Compress compresses the given byte slice.
func (c *LZ4Compressor) Compress(data []byte) ([]byte, error) {
	compressedData := make([]byte, lz4.CompressBlockBound(len(data)))
	n, err := lz4.CompressBlock(data, compressedData, nil)
	if err != nil {
		return nil, err
	}
	return compressedData[:n], nil
}

// Decompress decompress the given byte slice.
func (c *LZ4Compressor) Decompress(data []byte) ([]byte, error) {
	// FUTURE: For maximum safety, the original uncompressed size should be
	// transmitted as part of the message protocol. This prevents decompression
	// bombs and allocates the precise amount of memory needed.

	// Start with a reasonable buffer size (e.g., 3x compressed size).
	// This is a heuristic and might need tuning.
	bufferSize := len(data) * 3
	if bufferSize < 1024 { // Have a minimum buffer size
		bufferSize = 1024
	}

	decompressedData := make([]byte, bufferSize)

	n, err := lz4.UncompressBlock(data, decompressedData)
	if err != nil {
		// If the error indicates the buffer was too small, we could potentially
		// try again with a larger buffer, but this can be risky.
		// A better approach is to ensure the initial buffer is sufficient or
		// transmit the original size.
		return nil, fmt.Errorf("failed to decompress with lz4: %w", err)
	}

	return decompressedData[:n], nil
}
