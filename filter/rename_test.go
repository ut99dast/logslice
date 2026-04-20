package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewRenamer_Valid(t *testing.T) {
	r, err := NewRenamer([]string{"level=severity", "msg=message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Renamer")
	}
}

func TestNewRenamer_Invalid(t *testing.T) {
	cases := []struct {
		name  string
		exprs []string
	}{
		{"empty", []string{}},
		{"no equals", []string{"levelonly"}},
		{"same name", []string{"level=level"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRenamer(tc.exprs)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestRenamer_Apply(t *testing.T) {
	r, _ := NewRenamer([]string{"level=severity"})
	record := map[string]interface{}{"level": "info", "msg": "hello"}
	out := r.Apply(record)

	if _, ok := out["level"]; ok {
		t.Error("old field 'level' should not exist")
	}
	if v, ok := out["severity"]; !ok || v != "info" {
		t.Errorf("expected severity=info, got %v", out["severity"])
	}
	if v, ok := out["msg"]; !ok || v != "hello" {
		t.Errorf("expected msg=hello, got %v", v)
	}
}

func TestRenamer_Apply_MissingField(t *testing.T) {
	r, _ := NewRenamer([]string{"level=severity"})
	record := map[string]interface{}{"msg": "hello"}
	out := r.Apply(record)
	if _, ok := out["severity"]; ok {
		t.Error("severity should not be set when level is absent")
	}
}

func TestRunRename(t *testing.T) {
	input := `{"level":"info","msg":"started"}
{"level":"error","msg":"failed"}
`
	r, _ := NewRenamer([]string{"level=severity"})
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON, nil)
	err := RunRename(strings.NewReader(input), w, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Flush()
	out := buf.String()
	if strings.Contains(out, "\"level\"") {
		t.Error("output should not contain 'level' field")
	}
	if !strings.Contains(out, "\"severity\"") {
		t.Error("output should contain 'severity' field")
	}
}
