package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunWherePipeline_Basic(t *testing.T) {
	input := `{"level":"error","latency":250}
{"level":"info","latency":10}
{"level":"error","latency":50}
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunWherePipeline(r, &buf, []string{"level eq error", "latency gt 100"}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], `"latency":250`) {
		t.Errorf("unexpected output: %s", lines[0])
	}
}

func TestRunWherePipeline_NoExprs(t *testing.T) {
	r := strings.NewReader("{}\n")
	var buf bytes.Buffer
	err := RunWherePipeline(r, &buf, []string{}, FormatJSON)
	if err == nil {
		t.Error("expected error for empty expressions")
	}
}

func TestRunWherePipeline_SkipsInvalidLines(t *testing.T) {
	input := "not-json\n{\"score\":5}\n{\"score\":15}\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunWherePipeline(r, &buf, []string{"score gt 10"}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestRunWherePipeline_AllMatch(t *testing.T) {
	input := "{\"x\":1}\n{\"x\":2}\n{\"x\":3}\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunWherePipeline(r, &buf, []string{"x gte 1"}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestRunWherePipeline_NoneMatch(t *testing.T) {
	input := "{\"x\":1}\n{\"x\":2}\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunWherePipeline(r, &buf, []string{"x gt 100"}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got: %s", buf.String())
	}
}
