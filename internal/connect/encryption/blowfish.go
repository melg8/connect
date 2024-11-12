// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package encryption

import (
	"errors"

	"golang.org/x/crypto/blowfish"
)

type BlowFishKey struct {
	key []byte
}

func NewBlowFishKey(key []byte) *BlowFishKey {
	return &BlowFishKey{
		key: key,
	}
}

func DefaultAuthKey() *BlowFishKey {
	return NewBlowFishKey([]byte{
		0x5F, 0x3B, 0x35, 0x2E,
		0x5D, 0x39, 0x34, 0x2D,
		0x33, 0x31, 0x3D, 0x3D,
		0x2D, 0x25, 0x78, 0x54,
		0x21, 0x5E, 0x5B, 0x24,
		0x00})
}

func (b *BlowFishKey) Decrypt(encrypted []byte) ([]byte, error) {
	if b == nil || b.key == nil {
		return nil, errors.New("BlowFishKey or key is nil")
	}

	len := len(encrypted)
	if len == 0 {
		return nil, errors.New("encrypted data is empty")
	}

	if len%8 != 0 {
		return nil, errors.New("encrypted data is not a multiple of 8")
	}

	cipher, err := blowfish.NewCipher(b.key)
	if err != nil {
		return nil, errors.New("failed to initialize blowfish")
	}

	decrypted := make([]byte, len)
	count := len / 8

	for i := 0; i < count; i++ {
		cipher.Decrypt(decrypted[i*8:], encrypted[i*8:])
	}

	return decrypted, nil
}
