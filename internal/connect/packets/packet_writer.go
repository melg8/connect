// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"bytes"
	"encoding/binary"
)

type PacketWriter struct {
	bytes.Buffer
}

func NewPacketWriter() *PacketWriter {
	return &PacketWriter{}
}

func (b *PacketWriter) WriteInt64(value int64) error {
	return binary.Write(b, binary.LittleEndian, value)
}

func (b *PacketWriter) WriteInt32(value int32) error {
	return binary.Write(b, binary.LittleEndian, value)
}

func (b *PacketWriter) WriteInt16(value int16) error {
	return binary.Write(b, binary.LittleEndian, value)
}

func (b *PacketWriter) WriteInt8(value int8) error {
	return binary.Write(b, binary.LittleEndian, value)
}

func (b *PacketWriter) WriteBytes(bytes []byte) error {
	_, err := b.Write(bytes)
	return err
}
