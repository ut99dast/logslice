package filter

import "testing"

func TestMultiFilter_AllMatch(t *testing.T) {
	mf, err := NewMultiFilter([]string{"level=error", "service=auth"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := map[string]interface{}{
		"level":   "error",
		"service": "auth",
	}
	if !mf.Match(record) {
		t.Error("expected all filters to match")
	}
}

func TestMultiFilter_PartialMatch(t *testing.T) {
	mf, _ := NewMultiFilter([]string{"level=error", "service=auth"})
	record := map[string]interface{}{
		"level":   "error",
		"service": "payments",
	}
	if mf.Match(record) {
		t.Error("expected partial match to fail")
	}
}

func TestMultiFilter_Empty(t *testing.T) {
	mf, _ := NewMultiFilter([]string{})
	if !mf.Empty() {
		t.Error("expected Empty() to be true")
	}
	record := map[string]interface{}{"level": "info"}
	if !mf.Match(record) {
		t.Error("empty multi-filter should match everything")
	}
}

func TestMultiFilter_InvalidExpr(t *testing.T) {
	_, err := NewMultiFilter([]string{"level=error", "badexpr"})
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
