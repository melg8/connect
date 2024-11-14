// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packets

import (
	"bytes"
	"encoding/binary"
	"errors"

	"golang.org/x/text/encoding/unicode"
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

func (r *PacketReader) ReadInt32() (int32, error) {
	var buf [4]byte
	n, err := r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != 4 {
		return 0, errors.New("error: PacketReader.ReadInt32 not enough bytes to read")
	}
	result := int32(buf[3])<<24 |
		(int32(buf[2]) << 16) |
		(int32(buf[1]) << 8) |
		int32(buf[0])
	return result, nil
}

func (r *PacketReader) ReadInt16() (int16, error) {
	var result int16
	err := binary.Read(r, binary.LittleEndian, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *PacketReader) ReadInt8() (int8, error) {
	var result int8
	err := binary.Read(r, binary.LittleEndian, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r *PacketReader) ReadStringFromUtf16Format() (string, error) {
	var data []byte
	for {
		first_byte, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		second_byte, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if first_byte == 0 && second_byte == 0 {
			break
		}
		data = append(data, first_byte, second_byte)
	}
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	decodedString, err := decoder.String(string(data))
	if err != nil {
		return "", nil
	}
	return decodedString, nil
}
