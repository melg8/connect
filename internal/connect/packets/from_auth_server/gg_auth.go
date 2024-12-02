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

func (p *GGAuthPacket) ToBytes() ([]byte, error) {
	writer := packet.NewWriter()
	if err := writer.WriteInt32(p.SessionID); err != nil {
		return nil, err
	}
	if err := writer.WriteInt32(p.Unknown); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func (p *GGAuthPacket) ToString() string {
	return "\nGGAuthPacket:" +
		"\n  SessionID: " + helpers.HexStringFromInt32(p.SessionID) +
		"\n  Unknown: " + helpers.HexStringFromInt32(p.Unknown)
}