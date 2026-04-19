package filter

import (
	"testing"
)

func TestNewFieldFilter_Valid(t *testing.T) {
	tests := []struct {
		expr     string
		field    string
		op       string
		value    string
	}{
		{"level=error", "level", "=", "error"},
		{"status!=200", "status", "!=", "200"},
		{"service = auth", "service", "=", "auth"},
	}
	for _, tt := range tests {
		f, err := NewFieldFilter(tt.expr)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.expr, err)
		}
		if f.Field != tt.field || f.Operator != tt.op || f.Value != tt.value {
			t.Errorf("parsed %q got field=%q op=%q value=%q", tt.expr, f.Field, f.Operator, f.Value)
		}
	}
}

func TestNewFieldFilter_Invalid(t *testing.T) {
	_, err := NewFieldFilter("nodelmiter")
	if err == nil {
		t.Error("expected error for expression without operator")
	}
}

func TestFieldFilter_Match(t *testing.T) {
	record := map[string]interface{}{
		"level":  "error",
		"status": "500",
	}

	f, _ := NewFieldFilter("level=error")
	if !f.Match(record) {
		t.Error("expected match for level=error")
	}

	f2, _ := NewFieldFilter("level!=error")
	if f2.Match(record) {
		t.Error("expected no match for level!=error")
	}

	f3, _ := NewFieldFilter("missing=value")
	if f3.Match(record) {
		t.Error("expected no match for missing field")
	}
}
