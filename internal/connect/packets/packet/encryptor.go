// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packet

import (
	"errors"
	"log"

	"github.com/melg8/connect/internal/connect/crypt"
)

type Serializable interface {
	ToBytes(*Writer) error
}

const (
	messagePrefixSize = 2
	crcAllignBy       = 4
	crcSize           = 4
	cipherAllignBy    = 8
)

type Encryptor struct {
	writer Writer
	cipher *crypt.BlowFishCipher
}

func NewEncryptor(writer Writer, cipher *crypt.BlowFishCipher) *Encryptor {
	return &Encryptor{
		writer: writer,
		cipher: cipher,
	}
}

func (e *Encryptor) writePadding(count int) error {
	for i := 0; i < count; i++ {
		if err := e.writer.WriteInt8(0); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encryptor) writePaddingAndChecksum() error {
	currentMessageSize := e.writer.Len() - messagePrefixSize

	log.Printf("Current message size: %d", currentMessageSize)

	sizeWithCrc := currentMessageSize + crcSize
	paddingNeeded := (crcAllignBy - (sizeWithCrc % crcAllignBy)) % crcAllignBy

	log.Printf("Size with padding: %d", sizeWithCrc+paddingNeeded)

	if err := e.writePadding(paddingNeeded); err != nil {
		return err
	}

	crc, err := crypt.Checksum(e.writer.Bytes()[messagePrefixSize:])
	if err != nil {
		return err
	}
	crcBytes := make([]byte, 4)
	crcBytes[3] = byte(crc)
	crcBytes[2] = byte(crc >> 8)
	crcBytes[1] = byte(crc >> 16)
	crcBytes[0] = byte(crc >> 24)
	return e.writer.WriteBytes(crcBytes)
}

func (e *Encryptor) allignAndEncrypt() error {
	if e.writer.Len() < 2 {
		return errors.New("not enough data to encrypt")
	}

	paddingSize := (cipherAllignBy - (e.writer.Len()-2)%cipherAllignBy) % cipherAllignBy

	if err := e.writePadding(paddingSize); err != nil {
		return err
	}

	if err := e.cipher.EncryptInplace(e.writer.Bytes()[2:]); err != nil {
		return err
	}
	return nil
}

func (e *Encryptor) writePacketSize() error {
	packetSizeBytes := e.writer.Bytes()[0:2]
	packetSize := int16(e.writer.Len())
	packetSizeBytes[0] = byte(packetSize)
	packetSizeBytes[1] = byte(packetSize >> 8)
	return nil
}

func (e *Encryptor) Write(data Serializable) error {
	// Reserve 2 bytes for future size value
	e.writer.WriteInt16(0)

	if err := data.ToBytes(&e.writer); err != nil {
		return err
	}

	if err := e.writePaddingAndChecksum(); err != nil {
		return err
	}

	if err := e.allignAndEncrypt(); err != nil {
		return err
	}

	if err := e.writePacketSize(); err != nil {
		return err
	}

	return nil
}

func (e *Encryptor) Bytes() []byte {
	return e.writer.Bytes()
}
