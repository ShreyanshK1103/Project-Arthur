package builder

import (
	"os"
	"strings"
	"testing"
)

func TestRedact(t *testing.T) {
	e := &EnvVars{vars: map[string]string{
		"API_URL": "https://api.example.com",
		"TOKEN":   "abcd1234",
		"FLAG":    "1", // too short, must not be redacted
	}}

	got := e.Redact("calling https://api.example.com with abcd1234, flag=1")
	want := "calling *** with ***, flag=1"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRedactNilSafe(t *testing.T) {
	var e *EnvVars
	if got := e.Redact("hello"); got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestWriteEnvFile(t *testing.T) {
	e := &EnvVars{vars: map[string]string{"B": "2", "A": "x=y"}}

	path, cleanup, err := e.WriteEnvFile()
	if err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	if string(data) != "A=x=y\nB=2\n" {
		t.Fatalf("unexpected file contents %q", data)
	}

	info, _ := os.Stat(path)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}

	cleanup()
	if _, err := os.Stat(path); !strings.Contains(err.Error(), "no such file") {
		t.Fatal("expected file to be removed")
	}
}

func TestWriteEnvFileEmpty(t *testing.T) {
	path, cleanup, err := (&EnvVars{vars: map[string]string{}}).WriteEnvFile()
	defer cleanup()

	if err != nil || path != "" {
		t.Fatalf("expected empty path and no error, got %q, %v", path, err)
	}
}
