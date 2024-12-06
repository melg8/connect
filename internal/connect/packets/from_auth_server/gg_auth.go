// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package fromauthserver

import (
	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets/packet"
)

type GGAuthPacket struct {
	SessionID int32
	Unknown   int32
}

func NewGGAuthPacketFromBytes(data []byte) (*GGAuthPacket, error) {
	reader := packet.NewReader(data)
	sessionID, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	unknown, err := reader.ReadInt32()
	if err != nil {
		return nil, err
	}
	return &GGAuthPacket{
		SessionID: sessionID,
		Unknown:   unknown,
	}, nil
}

func (p *GGAuthPacket) ToBytes(writer *packet.Writer) error {
	if err := writer.WriteInt32(p.SessionID); err != nil {
		return err
	}
	if err := writer.WriteInt32(p.Unknown); err != nil {
		return err
	}
	return nil
}

func (p *GGAuthPacket) ToString() string {
	return "\nGGAuthPacket:" +
		"\n  SessionID: " + helpers.HexStringFromInt32(p.SessionID) +
		"\n  Unknown: " + helpers.HexStringFromInt32(p.Unknown)
}
