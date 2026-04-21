package filter

import (
	"testing"
)

func makeRec(kv ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func TestNewWhereFilter_Valid(t *testing.T) {
	_, err := NewWhereFilter("latency gt 100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewWhereFilter_Invalid(t *testing.T) {
	cases := []string{
		"",
		"latency gt",
		"latency badop 100",
		"too many parts here now",
	}
	for _, c := range cases {
		_, err := NewWhereFilter(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestWhereFilter_NumericOps(t *testing.T) {
	cases := []struct {
		expr    string
		val     interface{}
		want    bool
	}{
		{"latency gt 100", 200.0, true},
		{"latency gt 100", 50.0, false},
		{"latency gte 100", 100.0, true},
		{"latency lt 100", 50.0, true},
		{"latency lte 100", 100.0, true},
		{"latency eq 42", 42.0, true},
		{"latency ne 42", 99.0, true},
		{"latency ne 42", 42.0, false},
	}
	for _, c := range cases {
		f, err := NewWhereFilter(c.expr)
		if err != nil {
			t.Fatalf("%q: %v", c.expr, err)
		}
		rec := makeRec("latency", c.val)
		if got := f.Match(rec); got != c.want {
			t.Errorf("%q with %v: got %v, want %v", c.expr, c.val, got, c.want)
		}
	}
}

func TestWhereFilter_StringOps(t *testing.T) {
	f, _ := NewWhereFilter("status eq error")
	if !f.Match(makeRec("status", "error")) {
		t.Error("expected match")
	}
	if f.Match(makeRec("status", "ok")) {
		t.Error("expected no match")
	}
}

func TestWhereFilter_MissingField(t *testing.T) {
	f, _ := NewWhereFilter("missing eq value")
	if f.Match(makeRec("other", "x")) {
		t.Error("expected no match for missing field")
	}
}
