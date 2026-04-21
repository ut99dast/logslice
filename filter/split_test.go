package filter

import (
	"strings"
	"testing"
)

func TestNewSplitter_Valid(t *testing.T) {
	s, err := NewSplitter("tags", ",", "tag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.field != "tags" || s.delimiter != "," || s.outField != "tag" {
		t.Errorf("unexpected splitter state: %+v", s)
	}
}

func TestNewSplitter_DefaultOutField(t *testing.T) {
	s, err := NewSplitter("tags", ",", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.outField != "tags" {
		t.Errorf("expected outField=tags, got %q", s.outField)
	}
}

func TestNewSplitter_Invalid(t *testing.T) {
	if _, err := NewSplitter("", ",", ""); err == nil {
		t.Error("expected error for empty field")
	}
	if _, err := NewSplitter("tags", "", ""); err == nil {
		t.Error("expected error for empty delimiter")
	}
}

func TestSplitter_Apply_Basic(t *testing.T) {
	s, _ := NewSplitter("tags", ",", "tag")
	rec := map[string]interface{}{"tags": "go,rust,zig", "level": "info"}
	results, err := s.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 records, got %d", len(results))
	}
	expected := []string{"go", "rust", "zig"}
	for i, r := range results {
		if r["tag"] != expected[i] {
			t.Errorf("record %d: expected tag=%q, got %v", i, expected[i], r["tag"])
		}
		if _, ok := r["tags"]; ok {
			t.Errorf("record %d: original field 'tags' should be removed", i)
		}
		if r["level"] != "info" {
			t.Errorf("record %d: expected level=info, got %v", i, r["level"])
		}
	}
}

func TestSplitter_Apply_SkipsEmptyTokens(t *testing.T) {
	s, _ := NewSplitter("tags", ",", "tag")
	rec := map[string]interface{}{"tags": "go,,zig"}
	results, err := s.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 records, got %d", len(results))
	}
}

func TestSplitter_Apply_MissingField(t *testing.T) {
	s, _ := NewSplitter("tags", ",", "tag")
	rec := map[string]interface{}{"level": "info"}
	_, err := s.Apply(rec)
	if err == nil {
		t.Error("expected error for missing field")
	}
}

func TestSplitter_Apply_NonStringField(t *testing.T) {
	s, _ := NewSplitter("tags", ",", "tag")
	rec := map[string]interface{}{"tags": 42}
	_, err := s.Apply(rec)
	if err == nil {
		t.Error("expected error for non-string field")
	}
}

func TestRunSplit_Basic(t *testing.T) {
	input := `{"tags":"a,b,c","host":"srv1"}
{"tags":"x,y","host":"srv2"}
`
	var out strings.Builder
	scanner := NewScanner(strings.NewReader(input))
	writer := NewWriter(&out, FormatJSON)
	if err := RunSplit(scanner, writer, "tags", ",", "tag"); err != nil {
		t.Fatalf("RunSplit error: %v", err)
	}
	writer.Flush()
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 output lines, got %d: %v", len(lines), lines)
	}
}
