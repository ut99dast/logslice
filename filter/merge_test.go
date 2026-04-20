package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewMerger_Valid(t *testing.T) {
	m, err := NewMerger(`{"env":"prod"}`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestNewMerger_Invalid(t *testing.T) {
	cases := []struct {
		name   string
		input  string
	}{
		{"empty string", ""},
		{"bad JSON", "{not json}"},
		{"empty object", "{}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewMerger(tc.input, false)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestMerger_NoOverwrite(t *testing.T) {
	m, _ := NewMerger(`{"env":"prod","version":"1.0"}`, false)
	rec := map[string]interface{}{"msg": "hello", "env": "dev"}
	out := m.Apply(rec)
	if out["env"] != "dev" {
		t.Errorf("expected env=dev (no overwrite), got %v", out["env"])
	}
	if out["version"] != "1.0" {
		t.Errorf("expected version=1.0, got %v", out["version"])
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
}

func TestMerger_WithOverwrite(t *testing.T) {
	m, _ := NewMerger(`{"env":"prod"}`, true)
	rec := map[string]interface{}{"msg": "hello", "env": "dev"}
	out := m.Apply(rec)
	if out["env"] != "prod" {
		t.Errorf("expected env=prod (overwrite), got %v", out["env"])
	}
}

func TestRunMerge(t *testing.T) {
	input := `{"msg":"line1","env":"dev"}
{"msg":"line2"}
not-json
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunMerge(r, &buf, `{"env":"prod","app":"logslice"}`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
	var rec1 map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &rec1); err != nil {
		t.Fatalf("parse line1: %v", err)
	}
	if rec1["env"] != "dev" {
		t.Errorf("line1 env: expected dev (no overwrite), got %v", rec1["env"])
	}
	if rec1["app"] != "logslice" {
		t.Errorf("line1 app: expected logslice, got %v", rec1["app"])
	}
	var rec2 map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &rec2); err != nil {
		t.Fatalf("parse line2: %v", err)
	}
	if rec2["env"] != "prod" {
		t.Errorf("line2 env: expected prod, got %v", rec2["env"])
	}
}
