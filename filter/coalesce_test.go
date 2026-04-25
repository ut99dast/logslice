package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewCoalescer_Valid(t *testing.T) {
	_, err := NewCoalescer([]string{"a", "b"}, "out")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCoalescer_Invalid(t *testing.T) {
	if _, err := NewCoalescer([]string{"a"}, "out"); err == nil {
		t.Fatal("expected error for single field")
	}
	if _, err := NewCoalescer([]string{"a", "b"}, ""); err == nil {
		t.Fatal("expected error for empty outField")
	}
}

func TestCoalescer_FirstNonEmpty(t *testing.T) {
	c, _ := NewCoalescer([]string{"a", "b", "c"}, "result")
	rec := map[string]interface{}{"a": "", "b": "hello", "c": "world"}
	out := c.Apply(rec)
	if out["result"] != "hello" {
		t.Errorf("expected 'hello', got %v", out["result"])
	}
}

func TestCoalescer_AllEmpty(t *testing.T) {
	c, _ := NewCoalescer([]string{"a", "b"}, "result")
	rec := map[string]interface{}{"a": "", "b": ""}
	out := c.Apply(rec)
	if out["result"] != "" {
		t.Errorf("expected empty string, got %v", out["result"])
	}
}

func TestCoalescer_MissingFields(t *testing.T) {
	c, _ := NewCoalescer([]string{"x", "y"}, "result")
	rec := map[string]interface{}{"z": "value"}
	out := c.Apply(rec)
	if out["result"] != "" {
		t.Errorf("expected empty string for missing fields, got %v", out["result"])
	}
}

func TestCoalescer_NonStringValue(t *testing.T) {
	c, _ := NewCoalescer([]string{"a", "b"}, "result")
	rec := map[string]interface{}{"a": "", "b": 42.0}
	out := c.Apply(rec)
	if out["result"] != 42.0 {
		t.Errorf("expected 42.0, got %v", out["result"])
	}
}

func TestRunCoalesce(t *testing.T) {
	input := `{"a":"","b":"found"}` + "\n" +
		`{"a":"first","b":"second"}` + "\n" +
		`not-json` + "\n"

	var buf bytes.Buffer
	err := RunCoalesce(strings.NewReader(input), &buf, []string{"a", "b"}, "out", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}

	var rec0 map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &rec0); err != nil {
		t.Fatalf("failed to parse line 0: %v", err)
	}
	if rec0["out"] != "found" {
		t.Errorf("line 0: expected 'found', got %v", rec0["out"])
	}

	var rec1 map[string]interface{}
	if err := json.Unmarshal([]byte(lines[1]), &rec1); err != nil {
		t.Fatalf("failed to parse line 1: %v", err)
	}
	if rec1["out"] != "first" {
		t.Errorf("line 1: expected 'first', got %v", rec1["out"])
	}
}
