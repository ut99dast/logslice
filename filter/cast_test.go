package filter

import (
	"testing"
)

func TestNewCaster_Valid(t *testing.T) {
	for _, tc := range []struct {
		field  string
		target CastType
	}{
		{"count", CastInt},
		{"ratio", CastFloat},
		{"active", CastBool},
		{"name", CastString},
	} {
		_, err := NewCaster(tc.field, tc.target)
		if err != nil {
			t.Errorf("expected no error for field=%q target=%q, got %v", tc.field, tc.target, err)
		}
	}
}

func TestNewCaster_Invalid(t *testing.T) {
	if _, err := NewCaster("", CastInt); err == nil {
		t.Error("expected error for empty field")
	}
	if _, err := NewCaster("x", CastType("bytes")); err == nil {
		t.Error("expected error for unknown cast type")
	}
}

func TestCaster_ToInt(t *testing.T) {
	c, _ := NewCaster("count", CastInt)
	rec := map[string]interface{}{"count": "42", "msg": "hello"}
	out, err := c.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["count"] != int64(42) {
		t.Errorf("expected int64(42), got %v (%T)", out["count"], out["count"])
	}
	if out["msg"] != "hello" {
		t.Error("other fields should be preserved")
	}
}

func TestCaster_ToFloat(t *testing.T) {
	c, _ := NewCaster("ratio", CastFloat)
	rec := map[string]interface{}{"ratio": "3.14"}
	out, err := c.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ratio"].(float64) < 3.13 || out["ratio"].(float64) > 3.15 {
		t.Errorf("unexpected float value: %v", out["ratio"])
	}
}

func TestCaster_ToBool(t *testing.T) {
	c, _ := NewCaster("active", CastBool)
	rec := map[string]interface{}{"active": "true"}
	out, err := c.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["active"] != true {
		t.Errorf("expected true, got %v", out["active"])
	}
}

func TestCaster_MissingField(t *testing.T) {
	c, _ := NewCaster("count", CastInt)
	rec := map[string]interface{}{"msg": "hello"}
	out, err := c.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["count"]; ok {
		t.Error("missing field should not be added")
	}
}

func TestCaster_BadValue(t *testing.T) {
	c, _ := NewCaster("count", CastInt)
	rec := map[string]interface{}{"count": "not-a-number"}
	_, err := c.Apply(rec)
	if err == nil {
		t.Error("expected error for invalid int conversion")
	}
}

func TestRunCast(t *testing.T) {
	lines := []string{
		`{"count":"5","msg":"ok"}`,
		`not json`,
		`{"count":"10","msg":"done"}`,
	}
	c, _ := NewCaster("count", CastInt)
	results, err := RunCast(lines, c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0]["count"] != int64(5) {
		t.Errorf("expected int64(5), got %v", results[0]["count"])
	}
	if results[1]["count"] != int64(10) {
		t.Errorf("expected int64(10), got %v", results[1]["count"])
	}
}
