package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

// sign builds the header value GitHub sends: "sha256=" + hex(HMAC-SHA256(secret, body))
func sign(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestValidGithubSignature(t *testing.T) {
	const secret = "test-secret"
	t.Setenv("GITHUB_WEBHOOK_SECRET", secret)

	body := []byte(`{"ref":"refs/heads/main"}`)
	good := sign(secret, body)
	hexPart := strings.TrimPrefix(good, "sha256=")

	cases := []struct {
		name   string
		body   []byte
		header string
		want   bool
	}{
		{"valid signature", body, good, true},
		{"tampered body", []byte(`{"ref":"refs/heads/evil"}`), good, false},
		{"signed with wrong secret", body, sign("other-secret", body), false},
		{"missing header", body, "", false},
		{"missing sha256= prefix", body, hexPart, false},
		{"wrong algorithm prefix", body, "sha1=" + hexPart, false},
		{"not hex", body, "sha256=zzzz", false},
		{"truncated signature", body, good[:len(good)-2], false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := validGithubSignature(c.body, c.header); got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

// If the secret is not configured, an attacker can compute a "valid" HMAC
// using an empty key. The check must refuse everything in that case.
func TestValidGithubSignatureRejectsWhenSecretUnset(t *testing.T) {
	t.Setenv("GITHUB_WEBHOOK_SECRET", "")

	body := []byte(`{"ref":"refs/heads/main"}`)

	if validGithubSignature(body, sign("", body)) {
		t.Fatal("signature accepted even though no webhook secret is configured")
	}
}