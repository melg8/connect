// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packet

import "github.com/melg8/connect/internal/connect/crypt"

type Serializable interface {
	ToBytes(*Writer) error
}

const (
	messagePrefixSize = 2
	alignBy           = 8
	crcSize           = 4
)

type Encryptor struct {
	writer Writer
	key    *crypt.BlowFishCipher
}

func NewEncryptor(writer Writer, key *crypt.BlowFishCipher) *Encryptor {
	return &Encryptor{
		writer: writer,
		key:    key,
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
	sizeWithCrc := currentMessageSize + crcSize
	paddingNeeded := alignBy - (sizeWithCrc % alignBy)

	if err := e.writePadding(paddingNeeded); err != nil {
		return err
	}

	crc, err := crypt.Checksum(e.writer.Bytes()[messagePrefixSize:])
	if err != nil {
		return err
	}
	return e.writer.WriteInt32(int32(crc))
}

func (e *Encryptor) writePacketSize() error {
	packetSizeBytes := e.writer.Bytes()[0:2]
	return NewWriterTo(packetSizeBytes).WriteInt16(int16(e.writer.Len() - 2))
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

	if err := e.writePacketSize(); err != nil {
		return err
	}

	// if err := e.key.Encrypt(e.writer.Bytes()[2:]); err != nil {
	// 	return err
	// }

	// // Calculate size of packet and insert to first 2 bytes
	// packetSize := e.writer.Len()
	// e.writer.Reset()
	// e.writer.WriteInt16(int16(packetSize))

	return nil
}
