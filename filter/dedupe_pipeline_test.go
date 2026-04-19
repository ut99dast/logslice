package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDedupe_UniqueLines(t *testing.T) {
	input := `{"id":"1","msg":"a"}
{"id":"2","msg":"b"}
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	stats, err := RunDedupe(r, w, DedupeOptions{Fields: []string{"id"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Matched != 2 {
		t.Errorf("expected 2 matched, got %d", stats.Matched)
	}
}

func TestRunDedupe_DuplicateLines(t *testing.T) {
	input := `{"id":"1","msg":"a"}
{"id":"1","msg":"b"}
{"id":"2","msg":"c"}
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	stats, err := RunDedupe(r, w, DedupeOptions{Fields: []string{"id"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Matched != 2 {
		t.Errorf("expected 2 unique, got %d", stats.Matched)
	}
	if stats.Valid != 3 {
		t.Errorf("expected 3 valid, got %d", stats.Valid)
	}
}

func TestRunDedupe_InvalidLines(t *testing.T) {
	input := `not-json
{"id":"1"}
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	stats, err := RunDedupe(r, w, DedupeOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Invalid != 1 {
		t.Errorf("expected 1 invalid, got %d", stats.Invalid)
	}
	if stats.Matched != 1 {
		t.Errorf("expected 1 matched, got %d", stats.Matched)
	}
}
