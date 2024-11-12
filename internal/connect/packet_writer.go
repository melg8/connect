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

func (b *PacketWriter) WriteUInt64(value uint64) error {
	return binary.Write(b, binary.LittleEndian, value)
}

func (b *PacketWriter) WriteBytes(bytes []byte) error {
	_, err := b.Write(bytes)
	return err
}
