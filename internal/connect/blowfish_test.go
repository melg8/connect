// Test for blowfish

package connect

import (
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
	if err.Error() != "encrypted data is not a multiple of 8" {
		t.Fatal("Error message should be 'encrypted data is not a multiple of 8', got: ",
			err.Error())
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
