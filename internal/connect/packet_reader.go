// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type PacketReader struct {
	*bytes.Reader
}

func NewPacketReader(buffer []byte) *PacketReader {
	return &PacketReader{bytes.NewReader(buffer)}
}

func (r *PacketReader) ReadBytes(number int) ([]byte, error) {
	buffer := make([]byte, number)
	n, err := r.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n < number {
		return nil, errors.New("error: PacketReader.ReadBytes not enough bytes to read")
	}
	return buffer, nil
}

func (r *PacketReader) ReadUInt64() (uint64, error) {
	buffer := make([]byte, 8)
	n, err := r.Read(buffer)
	if err != nil {
		return 0, err
	}
	if n < 8 {
		return 0, errors.New("error: PacketReader.ReadUInt64 not enough bytes to read")
	}

	buf := bytes.NewBuffer(buffer)
	var result uint64
	binary.Read(buf, binary.LittleEndian, &result)
	return result, nil
}
