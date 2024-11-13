// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packets

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"
)

func OnlyPartialPacket(size int) []byte {
	result := make([]byte, size)
	return result
}

func OnlySessionIdPacket() []byte {
	packet := "ca0e617cffff"
	data, err := hex.DecodeString(packet)
	if err != nil {
		panic(err)
	}
	return data
}

func InitPacketData() []byte {
	packet := "ca0e617c21c6000040dbaa8d4c6f4b4ab33a538ff2977f24b1a08a2799b8696b2834efb8dfdf75807dfd14ef3051489fedf04712ba576139898c0a5de47c431ce407f5450092d747ff1e2e8294c0f00365f1b6a5d005767f9bda4ec694d43c7cc956ed4b2edc982d657c42a8793f3b1a6c4be631b97fd791a8a5adc9f9f4e8660dcba865b679b0124e95dd29fc9cc37720b6ad97f7e0bd07"
	data, err := hex.DecodeString(packet)
	if err != nil {
		panic(err)
	}
	return data
}

func ExpectedRsaPublicKey() []byte {
	rsaPublicKey := "40dbaa8d4c6f4b4ab33a538ff2977f24b1a08a2799b8696b2834efb8dfdf75807dfd14ef3051489fedf04712ba576139898c0a5de47c431ce407f5450092d747ff1e2e8294c0f00365f1b6a5d005767f9bda4ec694d43c7cc956ed4b2edc982d657c42a8793f3b1a6c4be631b97fd791a8a5adc9f9f4e8660dcba865b679b012"
	data, err := hex.DecodeString(rsaPublicKey)
	if err != nil {
		panic(err)
	}
	return data
}

func convertBytesToInt32BigEndian(bytes []byte) int32 {
	return int32(binary.BigEndian.Uint32(bytes))
}

func ExpectedGameGuard3() int32 {
	return convertBytesToInt32BigEndian([]byte{0x97, 0xad, 0xb6, 0x20})
}

func TestInitPacketEncodingAndDecoding(t *testing.T) {
	packetBin := InitPacketData()
	init_packet, err := NewInitPacketFromBytes(packetBin)
	if err != nil {
		panic(err)
	}
	if init_packet.sessionId != int32(0x7c610eca) {
		t.Error("wrong session id")
	}

	if init_packet.protocolVersion != int32(0x0000c621) {
		t.Error("wrong protocol version")
	}

	expectedRsaPublicKey := ExpectedRsaPublicKey()
	if !bytes.Equal(init_packet.rsaPublicKey, expectedRsaPublicKey) {
		t.Error("wrong rsa public key")
	}

	if init_packet.gameGuard1 != int32(0x29dd954e) {
		t.Error("wrong game guard part 1")
	}

	if init_packet.gameGuard2 != int32(0x77c39cfc) {
		t.Error("wrong game guard part 2")
	}

	if init_packet.gameGuard3 != ExpectedGameGuard3() {
		t.Error("wrong game guard part 3")
	}

	if init_packet.gameGuard4 != int32(0x07bde0f7) {
		t.Error("wrong game guard part 4")
	}

	if init_packet.blowfishKey != nil {
		t.Error("expected no blowfish key, but got one")
	}

	encodedPacket := init_packet.NewInitPacket()

	if !bytes.Equal(encodedPacket, packetBin) {
		t.Error("encoded packet is not equal to initial packet bin form")
	}

	blowFishKey := make([]byte, 21)
	init_packet.blowfishKey = &blowFishKey

	stringRepresentation := init_packet.ToString()
	t.Log(stringRepresentation)
	if !strings.Contains(stringRepresentation, "blowfishKey") {
		t.Error("expected blowfish key in string representation")
	}

	encodedPacket1 := init_packet.NewInitPacket()
	if len(encodedPacket)+21 != len(encodedPacket1) {
		t.Error("incorrect length after blowfish key addition")
	}
}

func TestInitPacketDecodingErrorOnEmptyPacket(t *testing.T) {
	emptyPacket := []byte{}
	_, err := NewInitPacketFromBytes(emptyPacket)
	if err == nil {
		t.Error("expected error on empty packet")
	}
}

func TestInitPacketDecodingErrorOnOnlyPartialPacket(t *testing.T) {
	normalPacketSize := 4*6 + 128
	for i := 0; i < normalPacketSize; i += 4 {
		onlySessionIdPacket := OnlyPartialPacket(i + 1)
		_, err := NewInitPacketFromBytes(onlySessionIdPacket)
		if err == nil {
			t.Errorf("expected error on partial packet size %d", i+1)
		}
	}
}

func TestInitPacketDecodingOkayWithOptionalBlowfishKeyPresent(t *testing.T) {
	normalPacketSize := 4*6 + 128
	blowFishKeySise := 21
	onlySessionIdPacket := OnlyPartialPacket(normalPacketSize + blowFishKeySise)
	_, err := NewInitPacketFromBytes(onlySessionIdPacket)
	if err != nil {
		t.Error("can't decode with blowfish key")
	}
}
