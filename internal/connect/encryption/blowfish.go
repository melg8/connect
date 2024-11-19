// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package encryption

import (
	"errors"
	"fmt"

	//nolint:staticcheck // Required for compatibility with legacy game protocol
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

	lenEncrypted := len(encrypted)
	if lenEncrypted == 0 {
		return nil, errors.New("encrypted data is empty")
	}

	cipher, err := blowfish.NewCipher(b.key)
	if err != nil {
		return nil, errors.New("failed to initialize blowfish")
	}

	blockSize := cipher.BlockSize()
	if lenEncrypted%blockSize != 0 {
		return nil, fmt.Errorf("encrypted data length must be a multiple of %d, got %d", blockSize, lenEncrypted)
	}

	decrypted := make([]byte, lenEncrypted)
	count := lenEncrypted / blockSize

	for i := 0; i < count; i++ {
		cipher.Decrypt(decrypted[i*blockSize:], encrypted[i*blockSize:])
	}

	return decrypted, nil
}

func (b *BlowFishKey) Encrypt(data []byte) ([]byte, error) {
	if b == nil || b.key == nil {
		return nil, errors.New("BlowFishKey or key is nil")
	}

	lenData := len(data)
	if lenData == 0 {
		return nil, errors.New("data is empty")
	}

	cipher, err := blowfish.NewCipher(b.key)
	if err != nil {
		return nil, errors.New("failed to initialize blowfish")
	}

	blockSize := cipher.BlockSize()
	if lenData%blockSize != 0 {
		return nil, fmt.Errorf("data length must be a multiple of %d, got %d", blockSize, lenData)
	}

	encrypted := make([]byte, lenData)
	count := lenData / blockSize

	for i := 0; i < count; i++ {
		cipher.Encrypt(encrypted[i*blockSize:], data[i*blockSize:])
	}

	return encrypted, nil
}
