// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package encryption

import (
	"strings"
	"testing"
)

func TestBlowFishEmptyData(t *testing.T) {
	encryptedData := []byte{}
	_, err := DefaultAuthKey().Decrypt(encryptedData)
	if err == nil {
		t.Fatal("Decryption should have failed on empty data")
	}
	if err.Error() != "encrypted data is empty" {
		t.Fatal("Error message should be 'encrypted data is empty', got: ",
			err.Error())
	}
}

func TestBlowFishDataNotMultipleOf8(t *testing.T) {
	encryptedData := []byte{1, 2, 3, 4, 5, 6, 7}
	_, err := DefaultAuthKey().Decrypt(encryptedData)
	if err == nil {
		t.Fatal("Decryption should have failed on data not multiple of 8")
	}
	expectedPrefix := "encrypted data length must be a multiple of 8"
	if !strings.HasPrefix(err.Error(), expectedPrefix) {
		t.Fatalf("Error message should start with '%s', got: %s", expectedPrefix, err.Error())
	}
}

func TestBlowFishEmptyKey(t *testing.T) {
	encryptedData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	_, err := NewBlowFishKey([]byte{}).Decrypt(encryptedData)
	if err == nil {
		t.Fatal("Decryption should have failed on empty key")
	}
	if err.Error() != "failed to initialize blowfish" {
		t.Fatal("Error message should be 'failed to initialize blowfish', got: ",
			err.Error())
	}
}

// Test BlowFish nil key.
func TestBlowFishNilKey(t *testing.T) {
	encryptedData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	_, err := NewBlowFishKey(nil).Decrypt(encryptedData)
	if err == nil {
		t.Fatal("Decryption should have failed on nil key")
	}
	if err.Error() != "BlowFishKey or key is nil" {
		t.Fatal("Error message should be 'BlowFishKey or key is nil', got: ",
			err.Error())
	}
}

func TestBlowFish(t *testing.T) {
	encryptedData := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	decryptedData, err := DefaultAuthKey().Decrypt(encryptedData)
	if err != nil {
		t.Fatal("Decryption failed: ", err)
	}
	if len(decryptedData) != 8 {
		t.Fatal("Decrypted data should be 8 bytes, got: ", len(decryptedData))
	}
}

func TestBlowFishEncryptEmptyData(t *testing.T) {
	data := []byte{}
	_, err := DefaultAuthKey().Encrypt(data)
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
	_, err := DefaultAuthKey().Encrypt(data)
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
	_, err := NewBlowFishKey([]byte{}).Encrypt(data)
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
	_, err := NewBlowFishKey(nil).Encrypt(data)
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
	key := DefaultAuthKey()

	// Test encryption
	encryptedData, err := key.Encrypt(originalData)
	if err != nil {
		t.Fatal("Encryption failed: ", err)
	}
	if len(encryptedData) != 8 {
		t.Fatal("Encrypted data should be 8 bytes, got: ", len(encryptedData))
	}

	// Test decryption of the encrypted data
	decryptedData, err := key.Decrypt(encryptedData)
	if err != nil {
		t.Fatal("Decryption failed: ", err)
	}
	if len(decryptedData) != 8 {
		t.Fatal("Decrypted data should be 8 bytes, got: ", len(decryptedData))
	}

	// Verify the decrypted data matches the original
	for i := 0; i < len(originalData); i++ {
		if originalData[i] != decryptedData[i] {
			t.Fatalf("Decrypted data differs from original at position %d: expected %d, got %d",
				i, originalData[i], decryptedData[i])
		}
	}
}
