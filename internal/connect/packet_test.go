// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"bytes"
	"testing"
)

func TestPacketWriterAndReader(t *testing.T) {
	writer := NewPacketWriter()
	int64Value := int64(1234567890123)
	err := writer.WriteInt64(int64Value)
	if err != nil {
		panic(err)
	}
	keyData := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	err = writer.WriteBytes(keyData)
	if err != nil {
		panic(err)
	}

	reader := NewPacketReader(writer.Bytes())

	gotInt64Value, err := reader.ReadInt64()
	if err != nil {
		panic(err)
	}
	if gotInt64Value != int64Value {
		t.Errorf("Got different UInt64 value: %d != %d", gotInt64Value, int64Value)
	}

	gotKeyData, err := reader.ReadBytes(len(keyData))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(gotKeyData, keyData) {
		t.Errorf("Got different key data: %s != %s", gotKeyData, keyData)
	}
}

func TestPacketReaderReadUInt64Error(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadInt64()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPacketReaderReadBytesError(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadBytes(1)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPacketReaderReadBytesNotEnoughBytesError(t *testing.T) {
	reader := NewPacketReader([]byte{1, 2, 3})
	_, err := reader.ReadBytes(4)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
