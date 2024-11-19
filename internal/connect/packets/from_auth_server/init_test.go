// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package from_auth_server

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

func OnlySessionIDPacket() []byte {
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
	RsaPublicKey := "40dbaa8d4c6f4b4ab33a538ff2977f24b1a08a2799b8696b2834efb8dfdf75807dfd14ef3051489fedf04712ba576139898c0a5de47c431ce407f5450092d747ff1e2e8294c0f00365f1b6a5d005767f9bda4ec694d43c7cc956ed4b2edc982d657c42a8793f3b1a6c4be631b97fd791a8a5adc9f9f4e8660dcba865b679b012"
	data, err := hex.DecodeString(RsaPublicKey)
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
		t.Fatal(err)
	}
	if init_packet.SessionID != int32(0x7c610eca) {
		t.Error("wrong session id")
	}

	if init_packet.ProtocolVersion != int32(0x0000c621) {
		t.Error("wrong protocol version")
	}

	expectedRsaPublicKey := ExpectedRsaPublicKey()
	if !bytes.Equal(init_packet.RsaPublicKey, expectedRsaPublicKey) {
		t.Error("wrong rsa public key")
	}

	if init_packet.GameGuard1 != int32(0x29dd954e) {
		t.Error("wrong game guard part 1")
	}

	if init_packet.GameGuard2 != int32(0x77c39cfc) {
		t.Error("wrong game guard part 2")
	}

	if init_packet.GameGuard3 != ExpectedGameGuard3() {
		t.Error("wrong game guard part 3")
	}

	if init_packet.GameGuard4 != int32(0x07bde0f7) {
		t.Error("wrong game guard part 4")
	}

	if init_packet.BlowfishKey != nil {
		t.Error("blowfish key should be nil")
	}

	// Test encoding
	encoded, err := init_packet.NewInitPacket()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, packetBin) {
		t.Error("wrong encoded packet")
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
	partialPacket := OnlySessionIDPacket()
	_, err := NewInitPacketFromBytes(partialPacket)
	if err == nil {
		t.Error("expected error on partial packet")
	}
}

func TestInitPacketDecodingOkayWithOptionalBlowfishKeyPresent(t *testing.T) {
	packetBin := InitPacketData()
	init_packet, err := NewInitPacketFromBytes(packetBin)
	if err != nil {
		t.Fatal(err)
	}
	blowFishKey := make([]byte, 21)
	init_packet.BlowfishKey = &blowFishKey

	encoded, err := init_packet.NewInitPacket()
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded) != len(packetBin)+21 {
		t.Error("incorrect length after blowfish key addition")
	}

	stringRepresentation := init_packet.ToString()
	if strings.Contains(stringRepresentation, "BlowfishKey: nil") {
		t.Error("expected BlowfishKey to be non-nil")
	}
}
