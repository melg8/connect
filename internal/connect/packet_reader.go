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

func (r *PacketReader) ReadInt64() (int64, error) {
	var result int64
	err := binary.Read(r, binary.LittleEndian, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
