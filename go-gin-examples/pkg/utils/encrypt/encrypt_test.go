package encrypt

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "changeme"
	encrypt, _ := Encrypt(password)

	if !ValidatePassword(encrypt, password) {
		t.Errorf("expected verification success, but failed")
	}
}
