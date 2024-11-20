// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package toauthserver

import (
	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets"
)

type RequestGGAuth struct {
	SessionID int32
	Data1     int32
	Data2     int32
	Data3     int32
	Data4     int32
}

func NewRequestGGAuth(data []byte) (*RequestGGAuth, error) {
	var result RequestGGAuth

	reader := packets.NewPacketReader(data)
	var err error
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

func (p *RequestGGAuth) ToBytes() ([]byte, error) {
	writer := new(packets.PacketWriter)
	if err := writer.WriteInt32(p.SessionID); err != nil {
		return nil, err
	}
	if err := writer.WriteInt32(p.Data1); err != nil {
		return nil, err
	}
	if err := writer.WriteInt32(p.Data2); err != nil {
		return nil, err
	}
	if err := writer.WriteInt32(p.Data3); err != nil {
		return nil, err
	}
	if err := writer.WriteInt32(p.Data4); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func (p *RequestGGAuth) ToString() string {
	return "\nRequestGGAuth:" +
		"\n  SessionID: " + helpers.HexStringFromInt32(p.SessionID) +
		"\n  Data1: " + helpers.HexStringFromInt32(p.Data1) +
		"\n  Data2: " + helpers.HexStringFromInt32(p.Data2) +
		"\n  Data3: " + helpers.HexStringFromInt32(p.Data3) +
		"\n  Data4: " + helpers.HexStringFromInt32(p.Data4)
}
