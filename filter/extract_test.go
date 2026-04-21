package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewExtractor_Valid(t *testing.T) {
	ext, err := NewExtractor("level,msg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ext.fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(ext.fields))
	}
}

func TestNewExtractor_Invalid(t *testing.T) {
	cases := []string{"", "  ,  "}
	for _, c := range cases {
		_, err := NewExtractor(c)
		if err == nil {
			t.Errorf("expected error for input %q", c)
		}
	}
}

func TestExtractor_Apply_KeepsFields(t *testing.T) {
	ext, _ := NewExtractor("level,msg")
	record := map[string]interface{}{
		"level": "info",
		"msg":   "hello",
		"ts":    "2024-01-01",
	}
	out := ext.Apply(record)
	if _, ok := out["level"]; !ok {
		t.Error("expected 'level' in output")
	}
	if _, ok := out["msg"]; !ok {
		t.Error("expected 'msg' in output")
	}
	if _, ok := out["ts"]; ok {
		t.Error("'ts' should not be in output")
	}
}

func TestExtractor_Apply_MissingField(t *testing.T) {
	ext, _ := NewExtractor("level,missing")
	record := map[string]interface{}{"level": "warn"}
	out := ext.Apply(record)
	if len(out) != 1 {
		t.Errorf("expected 1 field, got %d", len(out))
	}
}

func TestRunExtract(t *testing.T) {
	input := `{"level":"info","msg":"started","ts":"2024-01-01T00:00:00Z"}
{"level":"error","msg":"failed","ts":"2024-01-01T00:01:00Z"}
`
	var buf bytes.Buffer
	w, err := NewWriter(&buf, FormatJSON)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if err := RunExtract(strings.NewReader(input), w, "level,msg"); err != nil {
		t.Fatalf("RunExtract: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "level") {
		t.Error("expected 'level' in output")
	}
	if strings.Contains(output, "ts") {
		t.Error("'ts' should not appear in output")
	}
}
