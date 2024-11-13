// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packets

import (
	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets"
)

type InitPacket struct {
	sessionId       int32
	protocolVersion int32
	rsaPublicKey    []byte
	gameGuard1      int32
	gameGuard2      int32
	gameGuard3      int32
	gameGuard4      int32
	blowfishKey     *[]byte
}

func NewInitPacketFromBytes(data []byte) (*InitPacket, error) {
	var result InitPacket
	var err error
	reader := packets.NewPacketReader(data)
	result.sessionId, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.protocolVersion, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.rsaPublicKey, err = reader.ReadBytes(128)
	if err != nil {
		return nil, err
	}

	result.gameGuard1, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.gameGuard2, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.gameGuard3, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	result.gameGuard4, err = reader.ReadInt32()
	if err != nil {
		return nil, err
	}

	var blowFishKey []byte
	blowFishKey, err = reader.ReadBytes(21)
	if err == nil {
		result.blowfishKey = &blowFishKey
	}

	return &result, nil
}

func (p *InitPacket) NewInitPacket() []byte {
	buffer := new(packets.PacketWriter)

	buffer.WriteInt32(p.sessionId)
	buffer.WriteInt32(p.protocolVersion)
	buffer.WriteBytes(p.rsaPublicKey)
	buffer.WriteInt32(p.gameGuard1)
	buffer.WriteInt32(p.gameGuard2)
	buffer.WriteInt32(p.gameGuard3)
	buffer.WriteInt32(p.gameGuard4)
	if p.blowfishKey != nil {
		buffer.WriteBytes(*p.blowfishKey)
	}
	return buffer.Bytes()
}

func (p *InitPacket) ToString() string {
	result := "\nInitPacket:\n" +
		"  sessionId: " + helpers.HexStringFromInt32(p.sessionId) + "\n" +
		"  protocolVersion: " + helpers.HexStringFromInt32(p.protocolVersion) + "\n" +
		"  rsaPublicKey: " + helpers.HexViewFrom(p.rsaPublicKey) + "\n" +
		"  gameGuard1: " + helpers.HexStringFromInt32(p.gameGuard1) + "\n" +
		"  gameGuard2: " + helpers.HexStringFromInt32(p.gameGuard2) + "\n" +
		"  gameGuard3: " + helpers.HexStringFromInt32(p.gameGuard3) + "\n" +
		"  gameGuard4: " + helpers.HexStringFromInt32(p.gameGuard4) + "\n"

	if p.blowfishKey != nil {
		result += "  blowfishKey: " + helpers.HexViewFrom(*p.blowfishKey)
	} else {
		result += "  blowfishKey: " + "nil"
	}
	return result
}
