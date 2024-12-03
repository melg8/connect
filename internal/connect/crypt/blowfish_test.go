// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package crypt

import (
	"strings"
	"testing"
)

func TestBlowFishEmptyData(t *testing.T) {
	data := []byte{}
	err := DefaultAuthKey().Decrypt(data)
	if err == nil {
		t.Fatal("Decryption should have failed on empty data")
	}
	if err.Error() != "encrypted data is empty" {
		t.Fatal("Error message should be 'encrypted data is empty', got: ",
			err.Error())
	}
}

func TestBlowFishDataNotMultipleOf8(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7}
	err := DefaultAuthKey().Decrypt(data)
	if err == nil {
		t.Fatal("Decryption should have failed on data not multiple of 8")
	}
	expectedPrefix := "encrypted data length must be a multiple of 8"
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Fatalf("Error message should start with '%s', got: %s", expectedPrefix, err.Error())
	}
}

func TestBlowFishEmptyKey(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	err := NewBlowFishKey([]byte{}).Decrypt(data)
	if err == nil {
		t.Fatal("Decryption should have failed on empty key")
	}
	if err.Error() != "failed to initialize blowfish" {
		t.Fatal("Error message should be 'failed to initialize blowfish', got: ",
			err.Error())
	}
}

func TestBlowFishNilKey(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	err := NewBlowFishKey(nil).Decrypt(data)
	if err == nil {
		t.Fatal("Decryption should have failed on nil key")
	}
	if err.Error() != "BlowFishKey or key is nil" {
		t.Fatal("Error message should be 'BlowFishKey or key is nil', got: ",
			err.Error())
	}
}

func TestBlowFish(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	err := DefaultAuthKey().Decrypt(data)
	if err != nil {
		t.Fatal("Decryption failed: ", err)
	}
	if len(data) != 8 {
		t.Fatal("Data should be 8 bytes, got: ", len(data))
	}
}

func TestBlowFishEncryptEmptyData(t *testing.T) {
	data := []byte{}
	err := DefaultAuthKey().Encrypt(data)
	if err == nil {
		t.Fatal("Encryption should have failed on empty data")
	}
	if err.Error() != "data is empty" {
		t.Fatal("Error message should be 'data is empty', got: ",
			err.Error())
	}
}

func TestBlowFishEncryptDataNotMultipleOf8(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7}
	err := DefaultAuthKey().Encrypt(data)
	if err == nil {
		t.Fatal("Encryption should have failed on data not multiple of 8")
	}
	expectedPrefix := "data length must be a multiple of 8"
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Fatalf("Error message should start with '%s', got: %s", expectedPrefix, err.Error())
	}
}

func TestBlowFishEncryptEmptyKey(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	err := NewBlowFishKey([]byte{}).Encrypt(data)
	if err == nil {
		t.Fatal("Encryption should have failed on empty key")
	}
	if err.Error() != "failed to initialize blowfish" {
		t.Fatal("Error message should be 'failed to initialize blowfish', got: ",
			err.Error())
	}
}

func TestBlowFishEncryptNilKey(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	err := NewBlowFishKey(nil).Encrypt(data)
	if err == nil {
		t.Fatal("Encryption should have failed on nil key")
	}
	if err.Error() != "BlowFishKey or key is nil" {
		t.Fatal("Error message should be 'BlowFishKey or key is nil', got: ",
			err.Error())
	}
}

func TestBlowFishEncryptDecryptCycle(t *testing.T) {
	originalData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	data := make([]byte, len(originalData))
	copy(data, originalData)
	key := DefaultAuthKey()

	// Test encryption
	err := key.Encrypt(data)
	if err != nil {
		t.Fatal("Encryption failed: ", err)
	}
	if len(data) != 8 {
		t.Fatal("Data should be 8 bytes, got: ", len(data))
	}

	// Verify data was actually changed by encryption
	different := false
	for i := 0; i < len(originalData); i++ {
		if originalData[i] != data[i] {
			different = true
			break
		}
	}
	if !different {
		t.Fatal("Encrypted data is identical to original data")
	}

	// Test decryption of the encrypted data
	err = key.Decrypt(data)
	if err != nil {
		t.Fatal("Decryption failed: ", err)
	}
	if len(data) != 8 {
		t.Fatal("Data should be 8 bytes, got: ", len(data))
	}

	// Verify the decrypted data matches the original
	for i := 0; i < len(originalData); i++ {
		if originalData[i] != data[i] {
			t.Fatalf("Decrypted data differs from original at position %d: expected %d, got %d",
				i, originalData[i], data[i])
		}
	}
}
