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
	err := DefaultAuthKey().DecryptInplace(data)
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
	err := DefaultAuthKey().DecryptInplace(data)
	if err == nil {
		t.Fatal("Decryption should have failed on data not multiple of 8")
	}
	expectedPrefix := "encrypted data length must be a multiple of 8"
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Fatalf("Error message should start with '%s', got: %s", expectedPrefix, err.Error())
	}
}

func TestBlowFishEmptyKey(t *testing.T) {
	cipher, err := NewBlowFishCipher([]byte{})
	if err == nil {
		t.Fatal("Creating cipher should have failed with empty key")
	}
	if cipher != nil {
		t.Fatal("Cipher should be nil when creation fails")
	}
}

func TestBlowFishNilKey(t *testing.T) {
	cipher, err := NewBlowFishCipher(nil)
	if err == nil {
		t.Fatal("Creating cipher should have failed with nil key")
	}
	if cipher != nil {
		t.Fatal("Cipher should be nil when creation fails")
	}
}

func TestBlowFish(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	err := DefaultAuthKey().DecryptInplace(data)
	if err != nil {
		t.Fatal("Decryption failed: ", err)
	}
	if len(data) != 8 {
		t.Fatal("Data should be 8 bytes, got: ", len(data))
	}
}

func TestBlowFishEncryptEmptyData(t *testing.T) {
	data := []byte{}
	err := DefaultAuthKey().EncryptInplace(data)
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
	err := DefaultAuthKey().EncryptInplace(data)
	if err == nil {
		t.Fatal("Encryption should have failed on data not multiple of 8")
	}
	expectedPrefix := "data length must be a multiple of 8"
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Fatalf("Error message should start with '%s', got: %s", expectedPrefix, err.Error())
	}
}

func TestBlowFishEncryptEmptyKey(t *testing.T) {
	cipher, err := NewBlowFishCipher([]byte{})
	if err == nil {
		t.Fatal("Creating cipher should have failed with empty key")
	}
	if cipher != nil {
		t.Fatal("Cipher should be nil when creation fails")
	}
}

func TestBlowFishEncryptNilKey(t *testing.T) {
	cipher, err := NewBlowFishCipher(nil)
	if err == nil {
		t.Fatal("Creating cipher should have failed with nil key")
	}
	if cipher != nil {
		t.Fatal("Cipher should be nil when creation fails")
	}
}

func TestBlowFishEncryptDecryptCycle(t *testing.T) {
	originalData := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	data := make([]byte, len(originalData))
	copy(data, originalData)
	key := DefaultAuthKey()

	// Test encryption
	err := key.EncryptInplace(data)
	if err != nil {
		t.Fatal("Encryption failed: ", err)
	}
	if len(data) != 16 {
		t.Fatal("Data should be 16 bytes, got: ", len(data))
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
	err = key.DecryptInplace(data)
	if err != nil {
		t.Fatal("Decryption failed: ", err)
	}
	if len(data) != 16 {
		t.Fatal("Data should be 16 bytes, got: ", len(data))
	}

	// Verify the decrypted data matches the original
	for i := 0; i < len(originalData); i++ {
		if originalData[i] != data[i] {
			t.Fatalf("Decrypted data differs from original at position %d: expected %d, got %d",
				i, originalData[i], data[i])
		}
	}
}
