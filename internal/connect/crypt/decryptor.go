// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package crypt

import (
	"fmt"

	"github.com/melg8/connect/internal/connect/packets/packet"
)

type Deserializable interface {
	FromBytes(*packet.Reader) error
}

type Decryptor struct {
	reader *packet.Reader
	cipher *BlowfishCipher
}

func NewDecryptor(reader *packet.Reader, cipher *BlowfishCipher) *Decryptor {
	return &Decryptor{
		reader: reader,
		cipher: cipher,
	}
}

func (d *Decryptor) Read(_ Deserializable) error {
	size, err := d.reader.ReadInt16()
	if err != nil {
		return fmt.Errorf("failed to read packet size: %w", err)
	}

	if size < messagePrefixSize {
		return fmt.Errorf("invalid packet size: %d", size)
	}

	return nil
}
