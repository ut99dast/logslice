package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewTruncator_Invalid(t *testing.T) {
	_, err := NewTruncator("", 10, "...")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
	_, err = NewTruncator("msg", 0, "...")
	if err == nil {
		t.Fatal("expected error for zero maxLen")
	}
	_, err = NewTruncator("msg", -5, "...")
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestTruncator_ShortValue(t *testing.T) {
	tr, _ := NewTruncator("msg", 20, "...")
	rec := Record{"msg": "hello", "level": "info"}
	out, err := tr.Apply(rec)
	if err != nil {
		t.Fatal(err)
	}
	if out["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", out["msg"])
	}
}

func TestTruncator_LongValue(t *testing.T) {
	tr, _ := NewTruncator("msg", 5, "...")
	rec := Record{"msg": "this is a long message", "level": "warn"}
	out, err := tr.Apply(rec)
	if err != nil {
		t.Fatal(err)
	}
	if out["msg"] != "this ..." {
		t.Errorf("unexpected value: %v", out["msg"])
	}
}

func TestTruncator_MissingField(t *testing.T) {
	tr, _ := NewTruncator("msg", 5, "...")
	rec := Record{"level": "error"}
	out, err := tr.Apply(rec)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["msg"]; ok {
		t.Error("field should not be added")
	}
}

func TestTruncator_NonStringField(t *testing.T) {
	tr, _ := NewTruncator("count", 3, "...")
	rec := Record{"count": float64(42)}
	out, err := tr.Apply(rec)
	if err != nil {
		t.Fatal(err)
	}
	if out["count"] != float64(42) {
		t.Error("non-string field should be unchanged")
	}
}

func TestRunTruncate(t *testing.T) {
	input := `{"msg":"hello world","level":"info"}
{"msg":"hi","level":"debug"}
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	if err := RunTruncate(r, &buf, "msg", 5, "...", FormatJSON); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var rec Record
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatal(err)
	}
	if rec["msg"] != "hello..." {
		t.Errorf("expected 'hello...', got %v", rec["msg"])
	}
	if err := json.Unmarshal([]byte(lines[1]), &rec); err != nil {
		t.Fatal(err)
	}
	if rec["msg"] != "hi" {
		t.Errorf("expected 'hi', got %v", rec["msg"])
	}
}
