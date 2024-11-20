package toauthserver

import (
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

func TestRequestGGAuth_ToBytes(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
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

		// Expected byte sequence
		expected := []byte{
			0x01, 0x00, 0x00, 0x00, // SessionID: 1
			0x02, 0x00, 0x00, 0x00, // Data1: 2
			0x03, 0x00, 0x00, 0x00, // Data2: 3
			0x04, 0x00, 0x00, 0x00, // Data3: 4
			0x05, 0x00, 0x00, 0x00, // Data4: 5
		}
		assert.Equal(t, expected, data)

		// Verify that we can read it back
		decoded, err := NewRequestGGAuth(data)
		assert.NoError(t, err)
		assert.Equal(t, req.SessionID, decoded.SessionID)
		assert.Equal(t, req.Data1, decoded.Data1)
		assert.Equal(t, req.Data2, decoded.Data2)
		assert.Equal(t, req.Data3, decoded.Data3)
		assert.Equal(t, req.Data4, decoded.Data4)
	})
}
