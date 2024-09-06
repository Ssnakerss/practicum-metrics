package compression

import (
	"bytes"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	// Test data
	data := []byte("Hello, world!")

	// Compress the data
	compressedData, err := Compress(data)
	if err != nil {
		t.Errorf("Compress failed: %v", err)
	}

	// Decompress the data
	decompressedData, err := Decompress(compressedData)
	if err != nil {
		t.Errorf("Decompress failed: %v", err)
	}

	// Compare the original and decompressed data
	if !bytes.Equal(data, decompressedData) {
		t.Errorf("Decompressed data does not match original data")
	}
}
