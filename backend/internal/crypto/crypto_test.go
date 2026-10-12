package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func setup(t *testing.T) {
	t.Helper()
	t.Setenv("ENV_ENCRIPTION_KEY", base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, 32)))
	if err := Init(); err != nil {
		t.Fatal(err)
	}
}

func TestRoundTrip(t *testing.T) {
	setup(t)

	enc, err := Encrypt("https://api.example.com")
	if err != nil {
		t.Fatal(err)
	}

	got, err := Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}

	if got != "https://api.example.com" {
		t.Fatalf("got %q", got)
	}
}

func TestEncryptionIsRandomised(t *testing.T) {
	setup(t)

	a, _ := Encrypt("same")
	b, _ := Encrypt("same")

	if bytes.Equal(a, b) {
		t.Fatal("two encryptions of the same value should differ")
	}
}

func TestTamperingDetected(t *testing.T) {
	setup(t)

	enc, _ := Encrypt("secret")
	enc[len(enc)-1] ^= 0xff

	if _, err := Decrypt(enc); err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestInitRejectsWrongKeyLength(t *testing.T) {
	t.Setenv("ENV_ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(make([]byte, 16)))

	if err := Init(); err == nil {
		t.Fatal("expected error for 16-byte key")
	}
}
