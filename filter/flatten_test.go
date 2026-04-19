package filter

import (
	"testing"
)

func TestNewFlattener_Invalid(t *testing.T) {
	_, err := NewFlattener(-1)
	if err == nil {
		t.Fatal("expected error for negative maxDepth")
	}
}

func TestFlattener_Flat(t *testing.T) {
	fl, _ := NewFlattener(0)
	rec := map[string]any{"a": "1", "b": "2"}
	out := fl.Flatten(rec)
	if out["a"] != "1" || out["b"] != "2" {
		t.Errorf("unexpected flat result: %v", out)
	}
}

func TestFlattener_Nested(t *testing.T) {
	fl, _ := NewFlattener(0)
	rec := map[string]any{
		"user": map[string]any{
			"name": "alice",
			"age":  30,
		},
		"level": "info",
	}
	out := fl.Flatten(rec)
	if out["user.name"] != "alice" {
		t.Errorf("expected user.name=alice, got %v", out["user.name"])
	}
	if out["user.age"] != 30 {
		t.Errorf("expected user.age=30, got %v", out["user.age"])
	}
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
}

func TestFlattener_MaxDepth(t *testing.T) {
	fl, _ := NewFlattener(1)
	rec := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	out := fl.Flatten(rec)
	// depth limit 1: a.b should remain as map, not expand further
	if _, ok := out["a.b.c"]; ok {
		t.Error("expected a.b.c not to be expanded beyond maxDepth=1")
	}
	if _, ok := out["a.b"]; !ok {
		t.Error("expected a.b to exist as a nested map value")
	}
}

func TestRunFlatten(t *testing.T) {
	lines := []string{
		`{"msg":"hello","ctx":{"id":"42"}}`,
		`not json`,
		`{"x":"1"}`,
	}
	results, err := RunFlatten(lines, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0]["ctx.id"] != "42" {
		t.Errorf("expected ctx.id=42, got %v", results[0]["ctx.id"])
	}
}
