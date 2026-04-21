package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestPivotResult_Add(t *testing.T) {
	p := NewPivotResult("level", "service")
	p.Add("info", "auth")
	p.Add("info", "auth")
	p.Add("error", "auth")
	p.Add("info", "api")

	if p.Cells["info"]["auth"] != 2 {
		t.Errorf("expected 2, got %d", p.Cells["info"]["auth"])
	}
	if p.Cells["error"]["auth"] != 1 {
		t.Errorf("expected 1, got %d", p.Cells["error"]["auth"])
	}
	if p.Cells["info"]["api"] != 1 {
		t.Errorf("expected 1, got %d", p.Cells["info"]["api"])
	}
}

func TestPivotResult_Finalize(t *testing.T) {
	p := NewPivotResult("level", "service")
	p.Add("warn", "db")
	p.Add("info", "api")
	p.Finalize()

	if len(p.RowKeys) != 2 {
		t.Fatalf("expected 2 row keys, got %d", len(p.RowKeys))
	}
	if p.RowKeys[0] != "info" || p.RowKeys[1] != "warn" {
		t.Errorf("unexpected row key order: %v", p.RowKeys)
	}
}

func TestRunPivot_Basic(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","service":"auth"}`,
		`{"level":"info","service":"api"}`,
		`{"level":"error","service":"auth"}`,
		`{"level":"info","service":"auth"}`,
	}, "\n")

	var out bytes.Buffer
	err := RunPivot(strings.NewReader(input), &out, "level", "service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "info") {
		t.Error("expected 'info' in output")
	}
	if !strings.Contains(result, "error") {
		t.Error("expected 'error' in output")
	}
	if !strings.Contains(result, "auth") {
		t.Error("expected 'auth' in output")
	}
}

func TestRunPivot_SkipsInvalidLines(t *testing.T) {
	input := strings.Join([]string{
		`not json`,
		`{"level":"info","service":"api"}`,
	}, "\n")

	var out bytes.Buffer
	if err := RunPivot(strings.NewReader(input), &out, "level", "service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "info") {
		t.Error("expected valid line to appear in output")
	}
}

func TestRunPivot_EmptyFields(t *testing.T) {
	err := RunPivot(strings.NewReader(""), &bytes.Buffer{}, "", "service")
	if err == nil {
		t.Error("expected error for empty rowField")
	}
}

func TestRunPivotFromArgs_Positional(t *testing.T) {
	input := `{"level":"info","svc":"db"}`
	var out bytes.Buffer
	err := RunPivotFromArgs(strings.NewReader(input), &out, []string{"level", "svc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "level") {
		t.Error("expected pivot header in output")
	}
}
