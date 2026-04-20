package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewMasker_Valid(t *testing.T) {
	m, err := NewMasker("password:full")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.field != "password" || m.mode != MaskFull {
		t.Errorf("unexpected masker state: %+v", m)
	}
}

func TestNewMasker_Invalid(t *testing.T) {
	cases := []string{
		"nocodon",
		":full",
		"field:unknown",
		"email:partial:abc",
	}
	for _, c := range cases {
		_, err := NewMasker(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestMasker_FullMode(t *testing.T) {
	m, _ := NewMasker("token:full")
	rec := map[string]interface{}{"token": "supersecret", "level": "info"}
	out := m.Apply(rec)
	if out["token"] != "***" {
		t.Errorf("expected masked token, got %v", out["token"])
	}
	if out["level"] != "info" {
		t.Errorf("other fields should be unchanged")
	}
}

func TestMasker_PartialMode(t *testing.T) {
	m, _ := NewMasker("email:partial:2")
	rec := map[string]interface{}{"email": "user@example.com"}
	out := m.Apply(rec)
	masked, _ := out["email"].(string)
	if !strings.HasPrefix(masked, "us") || !strings.HasSuffix(masked, "om") {
		t.Errorf("unexpected partial mask: %q", masked)
	}
}

func TestMasker_MissingField(t *testing.T) {
	m, _ := NewMasker("secret:full")
	rec := map[string]interface{}{"level": "warn"}
	out := m.Apply(rec)
	if _, ok := out["secret"]; ok {
		t.Error("field should not be added if missing")
	}
}

func TestMasker_NonStringField(t *testing.T) {
	m, _ := NewMasker("code:full")
	rec := map[string]interface{}{"code": 42}
	out := m.Apply(rec)
	if out["code"] != 42 {
		t.Error("non-string field should be left unchanged")
	}
}

func TestRunMask(t *testing.T) {
	input := `{"user":"alice","password":"s3cr3t"}
{"user":"bob","password":"hunter2"}
not-json
`
	m, _ := NewMasker("password:full")
	var buf bytes.Buffer
	err := RunMask(strings.NewReader(input), &buf, m)
	if err == nil {
		t.Error("expected error due to skipped invalid line")
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	for _, line := range lines {
		if strings.Contains(line, "s3cr3t") || strings.Contains(line, "hunter2") {
			t.Errorf("password not masked in line: %s", line)
		}
	}
}
