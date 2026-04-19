package filter

import (
	"bytes"
	"strings"
	"testing"
)

type alwaysMatch struct{}

func (alwaysMatch) Match(_ map[string]interface{}) bool { return true }

type neverMatch struct{}

func (neverMatch) Match(_ map[string]interface{}) bool { return false }

func TestPipelineWithStats_AllMatch(t *testing.T) {
	input := strings.NewReader("{\"level\":\"info\"}\n{\"level\":\"warn\"}\n")
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	ps := NewPipelineWithStats(input, alwaysMatch{}, w)
	if err := ps.Run(input); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if ps.Stats.Matched != 2 {
		t.Errorf("expected 2 matched, got %d", ps.Stats.Matched)
	}
	if ps.Stats.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", ps.Stats.Skipped)
	}
}

func TestPipelineWithStats_NoneMatch(t *testing.T) {
	input := strings.NewReader("{\"level\":\"debug\"}\n")
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	ps := NewPipelineWithStats(input, neverMatch{}, w)
	if err := ps.Run(input); err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if ps.Stats.Matched != 0 {
		t.Errorf("expected 0 matched, got %d", ps.Stats.Matched)
	}
	if ps.Stats.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", ps.Stats.Skipped)
	}
}

func TestPipelineWithStats_InvalidLines(t *testing.T) {
	input := strings.NewReader("not json\n{\"ok\":true}\n")
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	ps := NewPipelineWithStats(input, alwaysMatch{}, w)
	_ = ps.Run(input)
	if ps.Stats.Invalid < 1 {
		t.Errorf("expected at least 1 invalid, got %d", ps.Stats.Invalid)
	}
}
