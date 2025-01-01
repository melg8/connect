// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package fromauthserver

import (
	"testing"

	"github.com/melg8/connect/internal/connect/packets/packet"
)

func BenchmarkPlusPlus(b *testing.B) {
	for range b.N {
		value := 2 + 2
		if value != 4 {
			b.Fatal("value is not 4")
		}
	}
}

func dataForBenchmark() []byte {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	return data
}

func BenchmarkInitPacketParsing(b *testing.B) {
	data := dataForBenchmark()

	b.ResetTimer()

	for range b.N {
		packet, err := NewInitPacketFromBytes(data)
		if err != nil {
			panic(err)
		}
		if packet.SessionID != 50462976 {
			b.Fatal("wrong session id")
		}
	}
}

func BenchmarkInitPacket_ToBytes(b *testing.B) {
	rsaKey := make([]byte, 128)
	for i := range rsaKey {
		rsaKey[i] = byte(i % 256)
	}
	blowfishKey := make([]byte, 21)
	for i := range blowfishKey {
		blowfishKey[i] = byte(i % 256)
	}

	initPacket := &InitPacket{
		SessionID:       12345,
		ProtocolVersion: 1,
		RsaPublicKey:    rsaKey,
		GameGuard1:      1,
		GameGuard2:      2,
		GameGuard3:      3,
		GameGuard4:      4,
		BlowfishKey:     &blowfishKey,
	}

	b.ResetTimer()
	for range b.N {
		packetWriter := packet.NewWriter()
		err := initPacket.ToBytes(packetWriter)
		if err != nil {
			b.Fatal(err)
		}
	}
}
