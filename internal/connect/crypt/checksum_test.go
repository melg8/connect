package crypt

import (
	"bytes"
	"testing"
)

func TestChecksumConsistency(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected uint32
	}{
		{
			name:     "empty data",
			data:     []byte{},
			expected: 0,
		},
		{
			name:     "single byte",
			data:     []byte{1},
			expected: 0,
		},
		{
			name:     "small data",
			data:     []byte{1, 2, 3, 4, 5},
			expected: 0,
		},
		{
			name:     "multiple of 4",
			data:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
			expected: 67372044,
		},
		{
			name:     "large data",
			data:     bytes.Repeat([]byte{255}, 1024),
			expected: 0x00000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := Checksum(tt.data)

			if result != tt.expected {
				t.Errorf("Checksum(%q) = %d, want %d", tt.data, result, tt.expected)
			}
		})
	}
}
