package media

import (
	"bytes"
	"testing"
)

func TestNoopCompressor_PassThrough(t *testing.T) {
	compressor := NewNoopCompressor()
	payload := []byte("benchmark-payload")

	compressed, err := compressor.Compress(payload)
	if err != nil {
		t.Fatalf("Compress returned error: %v", err)
	}
	if !bytes.Equal(compressed, payload) {
		t.Fatalf("Compress payload mismatch: got %q want %q", compressed, payload)
	}

	decompressed, err := compressor.Decompress(payload)
	if err != nil {
		t.Fatalf("Decompress returned error: %v", err)
	}
	if !bytes.Equal(decompressed, payload) {
		t.Fatalf("Decompress payload mismatch: got %q want %q", decompressed, payload)
	}
}
