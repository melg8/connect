package fromauthserver

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGGAuthPacketFromBytes(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		want        *GGAuthPacket
		wantErr     bool
		description string
	}{
		{
			name:  "valid packet",
			input: []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
			want: &GGAuthPacket{
				SessionID: 1,
				Unknown:   2,
			},
			wantErr:     false,
			description: "should successfully parse valid input",
		},
		{
			name:        "incomplete data for SessionID",
			input:       []byte{0x01, 0x00, 0x00},
			want:        nil,
			wantErr:     true,
			description: "should return error when not enough bytes for SessionID",
		},
		{
			name:        "incomplete data for Unknown",
			input:       []byte{0x01, 0x00, 0x00, 0x00, 0x02},
			want:        nil,
			wantErr:     true,
			description: "should return error when not enough bytes for Unknown field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGGAuthPacketFromBytes(tt.input)
			if tt.wantErr {
				assert.Error(t, err, tt.description)
				return
			}
			assert.NoError(t, err, tt.description)
			assert.Equal(t, tt.want, got, tt.description)
		})
	}
}

func TestGGAuthPacket_ToBytes(t *testing.T) {
	tests := []struct {
		name        string
		packet      *GGAuthPacket
		want        []byte
		wantErr     bool
		description string
	}{
		{
			name: "valid conversion",
			packet: &GGAuthPacket{
				SessionID: 1,
				Unknown:   2,
			},
			want:        []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
			wantErr:     false,
			description: "should successfully convert to bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.packet.ToBytes()
			if tt.wantErr {
				assert.Error(t, err, tt.description)
				return
			}
			assert.NoError(t, err, tt.description)
			assert.True(t, bytes.Equal(tt.want, got), tt.description)
		})
	}
}

func TestGGAuthPacket_ToString(t *testing.T) {
	packet := &GGAuthPacket{
		SessionID: 1,
		Unknown:   2,
	}
	result := packet.ToString()
	assert.Contains(t, result, "GGAuthPacket")
	assert.Contains(t, result, "SessionID: 00000001")
	assert.Contains(t, result, "Unknown: 00000002")
}

func TestGGAuthPacket_RoundTrip(t *testing.T) {
	original := &GGAuthPacket{
		SessionID: 12345,
		Unknown:   67890,
	}

	// Convert to bytes
	bytes, err := original.ToBytes()
	assert.NoError(t, err, "Failed to convert to bytes")

	// Convert back to packet
	reconstructed, err := NewGGAuthPacketFromBytes(bytes)
	assert.NoError(t, err, "Failed to parse bytes")

	// Compare
	assert.Equal(t, original.SessionID, reconstructed.SessionID, "SessionID mismatch after round trip")
	assert.Equal(t, original.Unknown, reconstructed.Unknown, "Unknown field mismatch after round trip")
}
