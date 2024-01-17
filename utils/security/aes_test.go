package security

import (
	"testing"
)

func Test_generateKey(t *testing.T) {
	data := "hello"

	encrypted, err := GetAESEncrypted(data)
	if err != nil {
		t.Error(err)
	}
	decrypted, err := GetAESDecrypted(encrypted)
	if err != nil {
		t.Error(err)
	}
	if string(decrypted) != data {
		t.Error("decrypted data is not equal to original data")
	}
}
