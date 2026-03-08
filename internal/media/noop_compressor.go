package media

// NoopCompressor implements domain.Compressor by passing bytes through unchanged.
type NoopCompressor struct{}

// NewNoopCompressor creates a new NoopCompressor.
func NewNoopCompressor() *NoopCompressor {
	return &NoopCompressor{}
}

// Compress returns the input as-is.
func (c *NoopCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

// Decompress returns the input as-is.
func (c *NoopCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
