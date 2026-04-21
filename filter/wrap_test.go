package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewWrapper_Valid(t *testing.T) {
	w, err := NewWrapper("msg", "[", "]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil Wrapper")
	}
}

func TestNewWrapper_EmptyField(t *testing.T) {
	_, err := NewWrapper("", "[", "]")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestWrapper_Apply_Basic(t *testing.T) {
	w, _ := NewWrapper("msg", "[", "]")
	rec := Record{"msg": "hello", "level": "info"}
	out := w.Apply(rec)
	if out["msg"] != "[hello]" {
		t.Errorf("expected [hello], got %v", out["msg"])
	}
	if out["level"] != "info" {
		t.Errorf("level should be unchanged")
	}
}

func TestWrapper_Apply_MissingField(t *testing.T) {
	w, _ := NewWrapper("missing", "[", "]")
	rec := Record{"msg": "hello"}
	out := w.Apply(rec)
	if _, ok := out["missing"]; ok {
		t.Error("missing field should not be created")
	}
}

func TestWrapper_Apply_NonStringField(t *testing.T) {
	w, _ := NewWrapper("count", "(", ")")
	rec := Record{"count": float64(42)}
	out := w.Apply(rec)
	if out["count"] != float64(42) {
		t.Error("non-string field should be unchanged")
	}
}

func TestRunWrapPipeline_Basic(t *testing.T) {
	input := `{"msg":"hello","level":"info"}` + "\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunWrapPipeline(r, &buf, []string{"msg:[:]", "level:<<:>>"}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[hello]") {
		t.Errorf("expected [hello] in output, got: %s", out)
	}
	if !strings.Contains(out, "<<info>>") {
		t.Errorf("expected <<info>> in output, got: %s", out)
	}
}

func TestRunWrapPipeline_NoExprs(t *testing.T) {
	err := RunWrapPipeline(strings.NewReader(""), &bytes.Buffer{}, nil, FormatJSON)
	if err == nil {
		t.Fatal("expected error for empty expressions")
	}
}

func TestRunWrapPipeline_InvalidExpr(t *testing.T) {
	err := RunWrapPipeline(strings.NewReader(""), &bytes.Buffer{}, []string{"badexpr"}, FormatJSON)
	if err == nil {
		t.Fatal("expected error for invalid expression format")
	}
}

func TestRunWrapPipeline_SkipsInvalidLines(t *testing.T) {
	input := "not-json\n" + `{"msg":"ok"}` + "\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	err := RunWrapPipeline(r, &buf, []string{"msg:>:<"}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 output line, got %d", len(lines))
	}
}
