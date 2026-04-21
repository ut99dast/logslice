package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewFormatter_Valid(t *testing.T) {
	_, err := NewFormatter("msg", "hello {name}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewFormatter_Invalid(t *testing.T) {
	cases := []struct {
		field, tmpl string
	}{
		{"", "hello {name}"},
		{"msg", ""},
		{"msg", "hello {name"},
		{"msg", "hello {}"},
	}
	for _, c := range cases {
		_, err := NewFormatter(c.field, c.tmpl)
		if err == nil {
			t.Errorf("expected error for field=%q tmpl=%q", c.field, c.tmpl)
		}
	}
}

func TestFormatter_Apply_Basic(t *testing.T) {
	fmt_, _ := NewFormatter("greeting", "Hello, {name}! You are {age} years old.")
	rec := map[string]interface{}{"name": "Alice", "age": 30}
	out, err := fmt_.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["greeting"] != "Hello, Alice! You are 30 years old." {
		t.Errorf("unexpected greeting: %v", out["greeting"])
	}
	// original fields preserved
	if out["name"] != "Alice" {
		t.Errorf("name field lost")
	}
}

func TestFormatter_Apply_LiteralOnly(t *testing.T) {
	fmt_, _ := NewFormatter("tag", "static-value")
	rec := map[string]interface{}{"x": 1}
	out, err := fmt_.Apply(rec)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["tag"] != "static-value" {
		t.Errorf("unexpected tag: %v", out["tag"])
	}
}

func TestFormatter_Apply_MissingField(t *testing.T) {
	fmt_, _ := NewFormatter("result", "{missing}")
	rec := map[string]interface{}{"x": 1}
	_, err := fmt_.Apply(rec)
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestFormatter_DoesNotMutateOriginal(t *testing.T) {
	fmt_, _ := NewFormatter("label", "{level}-event")
	rec := map[string]interface{}{"level": "info"}
	out, _ := fmt_.Apply(rec)
	out["level"] = "changed"
	if rec["level"] != "info" {
		t.Error("original record was mutated")
	}
}

func TestRunFormat(t *testing.T) {
	input := `{"host":"web1","status":200}
{"host":"db1","status":500}
`
	var buf bytes.Buffer
	err := RunFormat(strings.NewReader(input), &buf, "summary", "{host} returned {status}")
	if err != nil {
		t.Fatalf("RunFormat error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var rec map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &rec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if rec["summary"] != "web1 returned 200" {
		t.Errorf("unexpected summary: %v", rec["summary"])
	}
}
