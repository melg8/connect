package toauthserver

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestGGAuth(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		testData := []byte{
			0x01, 0x00, 0x00, 0x00, // SessionID: 1
			0x02, 0x00, 0x00, 0x00, // Data1: 2
			0x03, 0x00, 0x00, 0x00, // Data2: 3
			0x04, 0x00, 0x00, 0x00, // Data3: 4
			0x05, 0x00, 0x00, 0x00, // Data4: 5
		}

		req, err := NewRequestGGAuth(testData)
		assert.NoError(t, err)
		assert.NotNil(t, req)
		assert.Equal(t, int32(1), req.SessionID)
		assert.Equal(t, int32(2), req.Data1)
		assert.Equal(t, int32(3), req.Data2)
		assert.Equal(t, int32(4), req.Data3)
		assert.Equal(t, int32(5), req.Data4)
	})

	t.Run("empty data", func(t *testing.T) {
		req, err := NewRequestGGAuth([]byte{})
		assert.Error(t, err)
		assert.Nil(t, req)
	})

	t.Run("error reading SessionID", func(t *testing.T) {
		testData := []byte{0x01, 0x00, 0x00} // Incomplete SessionID
		req, err := NewRequestGGAuth(testData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})

	t.Run("error reading Data1", func(t *testing.T) {
		testData := []byte{
			0x01, 0x00, 0x00, 0x00, // Complete SessionID
			0x02, 0x00, 0x00, // Incomplete Data1
		}
		req, err := NewRequestGGAuth(testData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})

	t.Run("error reading Data2", func(t *testing.T) {
		testData := []byte{
			0x01, 0x00, 0x00, 0x00, // Complete SessionID
			0x02, 0x00, 0x00, 0x00, // Complete Data1
			0x03, 0x00, 0x00, // Incomplete Data2
		}
		req, err := NewRequestGGAuth(testData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})

	t.Run("error reading Data3", func(t *testing.T) {
		testData := []byte{
			0x01, 0x00, 0x00, 0x00, // Complete SessionID
			0x02, 0x00, 0x00, 0x00, // Complete Data1
			0x03, 0x00, 0x00, 0x00, // Complete Data2
			0x04, 0x00, 0x00, // Incomplete Data3
		}
		req, err := NewRequestGGAuth(testData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})

	t.Run("error reading Data4", func(t *testing.T) {
		testData := []byte{
			0x01, 0x00, 0x00, 0x00, // Complete SessionID
			0x02, 0x00, 0x00, 0x00, // Complete Data1
			0x03, 0x00, 0x00, 0x00, // Complete Data2
			0x04, 0x00, 0x00, 0x00, // Complete Data3
			0x05, 0x00, 0x00, // Incomplete Data4
		}
		req, err := NewRequestGGAuth(testData)
		assert.Error(t, err)
		assert.Nil(t, req)
	})
}

type mockPacketWriter struct {
	writeCount int
	failAfter  int
	buffer     bytes.Buffer
}

func newMockPacketWriter(failAfter int) *mockPacketWriter {
	return &mockPacketWriter{failAfter: failAfter}
}

func (w *mockPacketWriter) Write(p []byte) (n int, err error) {
	w.writeCount++
	if w.writeCount > w.failAfter {
		return 0, assert.AnError
	}
	return w.buffer.Write(p)
}

func (w *mockPacketWriter) WriteInt32(value int32) error {
	return binary.Write(w, binary.LittleEndian, value)
}

func (w *mockPacketWriter) Bytes() []byte {
	return w.buffer.Bytes()
}

func TestRequestGGAuth_ToBytes(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		req := &RequestGGAuth{
			SessionID: 1,
			Data1:     2,
			Data2:     3,
			Data3:     4,
			Data4:     5,
		}

		data, err := req.ToBytes()
		assert.NoError(t, err)
		assert.NotNil(t, data)

		// Verify the bytes are correct
		expected := []byte{
			0x01, 0x00, 0x00, 0x00, // SessionID: 1
			0x02, 0x00, 0x00, 0x00, // Data1: 2
			0x03, 0x00, 0x00, 0x00, // Data2: 3
			0x04, 0x00, 0x00, 0x00, // Data3: 4
			0x05, 0x00, 0x00, 0x00, // Data4: 5
		}
		assert.Equal(t, expected, data)

		// Test roundtrip
		newReq, err := NewRequestGGAuth(data)
		assert.NoError(t, err)
		assert.Equal(t, req, newReq)
	})

	t.Run("error cases", func(t *testing.T) {
		testCases := []struct {
			name      string
			failAfter int // number of successful writes before failing
			expectErr bool
		}{
			{"error writing SessionID", 0, true},
			{"error writing Data1", 1, true},
			{"error writing Data2", 2, true},
			{"error writing Data3", 3, true},
			{"error writing Data4", 4, true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				req := &RequestGGAuth{
					SessionID: 1,
					Data1:     2,
					Data2:     3,
					Data3:     4,
					Data4:     5,
				}

				// Create a new mock writer that fails after N writes
				writer := newMockPacketWriter(tc.failAfter)

				// Test error handling
				data, err := req.toBytesWithWriter(writer)
				if tc.expectErr {
					assert.Error(t, err)
					assert.Nil(t, data)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, data)
				}
			})
		}
	})
}

func TestRequestGGAuth_ToString(t *testing.T) {
	req := &RequestGGAuth{
		SessionID: 1,
		Data1:     2,
		Data2:     3,
		Data3:     4,
		Data4:     5,
	}

	str := req.ToString()
	assert.Contains(t, str, "SessionID: 00000001")
	assert.Contains(t, str, "Data1: 00000002")
	assert.Contains(t, str, "Data2: 00000003")
	assert.Contains(t, str, "Data3: 00000004")
	assert.Contains(t, str, "Data4: 00000005")
}
