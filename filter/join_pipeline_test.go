package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunJoinPipeline_SkipsInvalidLines(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "1", "tag": "alpha"},
	})
	j, _ := NewJoiner("id", path, false)

	input := "not-json\n{\"id\":\"1\",\"x\":1}\n"
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	n, err := RunJoinPipeline(strings.NewReader(input), j, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1 written, got %d", n)
	}
	if !strings.Contains(buf.String(), "alpha") {
		t.Errorf("expected enriched field 'alpha' in output")
	}
}

func TestRunJoinPipeline_EmptyInput(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "1", "tag": "beta"},
	})
	j, _ := NewJoiner("id", path, false)

	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	n, err := RunJoinPipeline(strings.NewReader(""), j, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 written, got %d", n)
	}
}

func TestRunJoinPipeline_AllNoMatch(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "99", "tag": "gamma"},
	})
	j, _ := NewJoiner("id", path, false)

	input := "{\"id\":\"1\",\"msg\":\"hi\"}\n{\"id\":\"2\",\"msg\":\"bye\"}\n"
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	n, err := RunJoinPipeline(strings.NewReader(input), j, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 written (pass-through), got %d", n)
	}
	if strings.Contains(buf.String(), "gamma") {
		t.Errorf("unexpected lookup field in output when no match")
	}
}
