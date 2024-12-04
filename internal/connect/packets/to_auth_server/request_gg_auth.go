// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package toauthserver

import (
	"errors"

	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets/packet"
)

const packetID = 0x07

type RequestGGAuth struct {
	SessionID int32
	Data1     int32
	Data2     int32
	Data3     int32
	Data4     int32
}

func Int32FromLEndian(s [4]byte) int32 {
	return int32(s[0]) |
		int32(s[1])<<8 |
		int32(s[2])<<16 |
		int32(s[3])<<24
}

func NewDefaultRequestGGAuth(sessionID int32) *RequestGGAuth {
	return &RequestGGAuth{
		SessionID: sessionID,
		// 23 92 90 4D 18 30 B5 7C 96 61 41 47 05 07 96 FB
		Data1: Int32FromLEndian([4]byte{0x23, 0x92, 0x90, 0x4D}),
		Data2: Int32FromLEndian([4]byte{0x18, 0x30, 0xB5, 0x7C}),
		Data3: Int32FromLEndian([4]byte{0x96, 0x61, 0x41, 0x47}),
		Data4: Int32FromLEndian([4]byte{0x05, 0x07, 0x96, 0xFB}),
	}
}

func NewRequestGGAuthFrom(data []byte) (*RequestGGAuth, error) {
	var result RequestGGAuth

	reader := packet.NewReader(data)

	id, err := reader.ReadInt8()
	if err != nil {
		return nil, err
	}
	if id != packetID {
		return nil, errors.New("invalid packet id")
	}

	result.SessionID, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.Data1, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.Data2, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.Data3, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.Data4, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *RequestGGAuth) ToBytes(writer *packet.Writer) error {
	if err := writer.WriteInt8(packetID); err != nil {
		return err
	}
	if err := writer.WriteInt32(p.SessionID); err != nil {
		return err
	}
	if err := writer.WriteInt32(p.Data1); err != nil {
		return err
	}
	if err := writer.WriteInt32(p.Data2); err != nil {
		return err
	}
	if err := writer.WriteInt32(p.Data3); err != nil {
		return err
	}
	if err := writer.WriteInt32(p.Data4); err != nil {
		return err
	}
	return nil
}

func (p *RequestGGAuth) ToString() string {
	return "\nRequestGGAuth:" +
		"\n  SessionID: " + helpers.HexStringFromInt32(p.SessionID) +
		"\n  Data1: " + helpers.HexStringFromInt32(p.Data1) +
		"\n  Data2: " + helpers.HexStringFromInt32(p.Data2) +
		"\n  Data3: " + helpers.HexStringFromInt32(p.Data3) +
		"\n  Data4: " + helpers.HexStringFromInt32(p.Data4)
}
