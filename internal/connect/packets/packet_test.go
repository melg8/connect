// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packets

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

	int32Value := int32(123456789)
	err = writer.WriteInt32(int32Value)
	if err != nil {
		panic(err)
	}

	int16Value := int16(12345)
	err = writer.WriteInt16(int16Value)
	if err != nil {
		panic(err)
	}

	int8Value := int8(123)
	err = writer.WriteInt8(int8Value)
	if err != nil {
		panic(err)
	}

	stringValue := "some string"
	err = writer.WriteStringAsUtf16(stringValue)
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

	gotInt32Value, err := reader.ReadInt32()
	if err != nil {
		panic(err)
	}
	if gotInt32Value != int32Value {
		t.Errorf("Got different int32 value: %d != %d", gotInt32Value, int32Value)
	}

	gotInt16Value, err := reader.ReadInt16()
	if err != nil {
		panic(err)
	}
	if gotInt16Value != int16Value {
		t.Errorf("Got different int16 value: %d != %d", gotInt16Value, int16Value)
	}

	gotInt8Value, err := reader.ReadInt8()
	if err != nil {
		panic(err)
	}
	if gotInt8Value != int8Value {
		t.Errorf("Got different int8 value: %d != %d", gotInt8Value, int8Value)
	}

	gotStringValue, err := reader.ReadStringFromUtf16Format()
	if err != nil {
		panic(err)
	}
	if gotStringValue != stringValue {
		t.Errorf("Got different string value: %s != %s", gotStringValue, stringValue)
	}

	gotKeyData, err := reader.ReadBytes(len(keyData))
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(gotKeyData, keyData) {
		t.Errorf("Got different key data: %s != %s", gotKeyData, keyData)
	}
}

func TestUtf16StringToHexAndAscii(t *testing.T) {
	writer := NewPacketWriter()

	stringValue := "some string"
	err := writer.WriteStringAsUtf16(stringValue)
	if err != nil {
		panic(err)
	}

	data := writer.Bytes()

	reader := NewPacketReader(data)

	gotStringValue, err := reader.ReadStringFromUtf16Format()
	if err != nil {
		panic(err)
	}
	if gotStringValue != stringValue {
		t.Errorf("Got different string value: %s != %s", gotStringValue, stringValue)
	}
}

func TestPacketReaderReadInt64Error(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadInt64()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPacketReaderReadInt32Error(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadInt32()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPacketReaderReadInt16Error(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadInt16()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPacketReaderReadInt8Error(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadInt8()
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

func TestPacketReaderReadStringFromUtf16FormatError(t *testing.T) {
	reader := NewPacketReader([]byte{})
	_, err := reader.ReadStringFromUtf16Format()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestPacketReaderReadStringFromUtf16FormatNotEnoughBytesError(t *testing.T) {
	reader := NewPacketReader([]byte{0x22, 0x00, 0x33})
	_, err := reader.ReadStringFromUtf16Format()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
