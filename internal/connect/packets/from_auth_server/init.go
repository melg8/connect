// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package from_auth_server

import (
	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets"
)

type InitPacket struct {
	SessionId       int32
	ProtocolVersion int32
	RsaPublicKey    []byte
	GameGuard1      int32
	GameGuard2      int32
	GameGuard3      int32
	GameGuard4      int32
	BlowfishKey     *[]byte
}

func NewInitPacketFromBytes(data []byte) (*InitPacket, error) {
	var result InitPacket
	var err error
	reader := packets.NewPacketReader(data)
	result.SessionId, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.ProtocolVersion, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.RsaPublicKey, err = reader.ReadBytes(128)
	if err != nil {
		return nil, err
	}
	result.GameGuard1, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.GameGuard2, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.GameGuard3, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	result.GameGuard4, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	var blowFishKey []byte
	blowFishKey, err = reader.ReadBytes(21)
	if err == nil {
		result.BlowfishKey = &blowFishKey
	}
	return &result, nil
}

func (p *InitPacket) NewInitPacket() []byte {
	buffer := new(packets.PacketWriter)
	buffer.WriteInt32(p.SessionId)
	buffer.WriteInt32(p.ProtocolVersion)
	buffer.WriteBytes(p.RsaPublicKey)
	buffer.WriteInt32(p.GameGuard1)
	buffer.WriteInt32(p.GameGuard2)
	buffer.WriteInt32(p.GameGuard3)
	buffer.WriteInt32(p.GameGuard4)
	if p.BlowfishKey != nil {
		buffer.WriteBytes(*p.BlowfishKey)
	}
	return buffer.Bytes()
}

func (p *InitPacket) ToString() string {
	result := "\nInitPacket:" +
		"\n  SessionId: " + helpers.HexStringFromInt32(p.SessionId) +
		"\n  ProtocolVersion: " + helpers.HexStringFromInt32(p.ProtocolVersion) +
		"\n  RsaPublicKey: \n" + helpers.HexViewFromWithLineSplit(p.RsaPublicKey, 16, "    ") +
		"\n  GameGuard1: " + helpers.HexStringFromInt32(p.GameGuard1) +
		"\n  GameGuard2: " + helpers.HexStringFromInt32(p.GameGuard2) +
		"\n  GameGuard3: " + helpers.HexStringFromInt32(p.GameGuard3) +
		"\n  GameGuard4: " + helpers.HexStringFromInt32(p.GameGuard4)

	if p.BlowfishKey != nil {
		result += "\n  BlowfishKey: \n" + helpers.HexViewFromWithLineSplit(*p.BlowfishKey, 16, "    ")
	} else {
		result += "\n  BlowfishKey: " + "nil"
	}
	return result
}
