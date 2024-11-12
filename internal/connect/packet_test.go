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
	uInt64Value := uint64(1234567890123)
	err := writer.WriteUInt64(uInt64Value)
	if err != nil {
		panic(err)
	}
	keyData := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	err = writer.WriteBytes(keyData)
	if err != nil {
		panic(err)
	}

	reader := NewPacketReader(writer.Bytes())

	gotUInt64Value, err := reader.ReadUInt64()
	if err != nil {
		panic(err)
	}
	if gotUInt64Value != uInt64Value {
		t.Errorf("Got different UInt64 value: %d != %d", gotUInt64Value, uInt64Value)
	}

	gotKeyData, err := reader.ReadBytes(len(keyData))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(gotKeyData, keyData) {
		t.Errorf("Got different key data: %s != %s", gotKeyData, keyData)
	}
}
