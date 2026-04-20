package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewHighlighter_Valid(t *testing.T) {
	_, err := NewHighlighter("msg", "error", HighlightBracket)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewHighlighter_Invalid(t *testing.T) {
	cases := []struct {
		field, pattern string
		mode           HighlightMode
	}{
		{"", "error", HighlightBracket},
		{"msg", "", HighlightBracket},
		{"msg", "(", HighlightBracket},
		{"msg", "error", "neon"},
	}
	for _, c := range cases {
		_, err := NewHighlighter(c.field, c.pattern, c.mode)
		if err == nil {
			t.Errorf("expected error for (%q,%q,%q)", c.field, c.pattern, c.mode)
		}
	}
}

func TestHighlighter_Bracket(t *testing.T) {
	h, _ := NewHighlighter("msg", "error", HighlightBracket)
	rec := map[string]interface{}{"msg": "an error occurred"}
	out := h.Apply(rec)
	if got := out["msg"].(string); got != "an [[error]] occurred" {
		t.Errorf("got %q", got)
	}
}

func TestHighlighter_Upper(t *testing.T) {
	h, _ := NewHighlighter("msg", "error", HighlightUpper)
	rec := map[string]interface{}{"msg": "an error occurred"}
	out := h.Apply(rec)
	if got := out["msg"].(string); got != "an ERROR occurred" {
		t.Errorf("got %q", got)
	}
}

func TestHighlighter_Mark(t *testing.T) {
	h, _ := NewHighlighter("msg", "error", HighlightMark)
	rec := map[string]interface{}{"msg": "an error occurred"}
	out := h.Apply(rec)
	if got := out["msg"].(string); got != ">>>an error occurred" {
		t.Errorf("got %q", got)
	}
}

func TestHighlighter_MissingField(t *testing.T) {
	h, _ := NewHighlighter("msg", "error", HighlightBracket)
	rec := map[string]interface{}{"level": "info"}
	out := h.Apply(rec)
	if _, ok := out["msg"]; ok {
		t.Error("expected no msg field in output")
	}
}

func TestHighlighter_NoMatch(t *testing.T) {
	h, _ := NewHighlighter("msg", "error", HighlightBracket)
	rec := map[string]interface{}{"msg": "all good"}
	out := h.Apply(rec)
	if got := out["msg"].(string); got != "all good" {
		t.Errorf("expected unchanged value, got %q", got)
	}
}

func TestRunHighlight(t *testing.T) {
	input := `{"msg":"an error occurred","level":"error"}
{"msg":"all good","level":"info"}
not-json
`
	h, _ := NewHighlighter("msg", "error", HighlightBracket)
	var buf bytes.Buffer
	if err := RunHighlight(strings.NewReader(input), &buf, h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var rec map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if got := rec["msg"].(string); got != "an [[error]] occurred" {
		t.Errorf("got %q", got)
	}
}
