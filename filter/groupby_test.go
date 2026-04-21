package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestGroupByResult_Add(t *testing.T) {
	g := NewGroupByResult("level")
	g.Add(map[string]interface{}{"level": "info", "msg": "a"})
	g.Add(map[string]interface{}{"level": "info", "msg": "b"})
	g.Add(map[string]interface{}{"level": "error", "msg": "c"})
	g.Add(map[string]interface{}{"msg": "no level"})

	if g.Counts["info"] != 2 {
		t.Errorf("expected info count 2, got %d", g.Counts["info"])
	}
	if g.Counts["error"] != 1 {
		t.Errorf("expected error count 1, got %d", g.Counts["error"])
	}
	if g.Counts["(missing)"] != 1 {
		t.Errorf("expected missing count 1, got %d", g.Counts["(missing)"])
	}
}

func TestGroupByResult_Sorted(t *testing.T) {
	g := NewGroupByResult("level")
	g.Counts["warn"] = 1
	g.Counts["info"] = 5
	g.Counts["error"] = 3

	sorted := g.Sorted()
	if sorted[0] != "info" {
		t.Errorf("expected first key 'info', got %q", sorted[0])
	}
	if sorted[1] != "error" {
		t.Errorf("expected second key 'error', got %q", sorted[1])
	}
	if sorted[2] != "warn" {
		t.Errorf("expected third key 'warn', got %q", sorted[2])
	}
}

func TestRunGroupBy_Basic(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"start"}`,
		`{"level":"error","msg":"fail"}`,
		`{"level":"info","msg":"end"}`,
		`{"msg":"no level"}`,
	}, "\n")

	var buf bytes.Buffer
	err := RunGroupBy(strings.NewReader(input), &buf, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "info") {
		t.Errorf("expected 'info' in output, got: %s", out)
	}
	if !strings.Contains(out, "error") {
		t.Errorf("expected 'error' in output, got: %s", out)
	}
	if !strings.Contains(out, "(missing)") {
		t.Errorf("expected '(missing)' in output, got: %s", out)
	}
}

func TestRunGroupBy_EmptyField(t *testing.T) {
	var buf bytes.Buffer
	err := RunGroupBy(strings.NewReader(`{"level":"info"}`), &buf, "")
	if err == nil {
		t.Error("expected error for empty field, got nil")
	}
}

func TestRunGroupBy_SkipsInvalidLines(t *testing.T) {
	input := strings.Join([]string{
		`not json`,
		`{"level":"debug"}`,
	}, "\n")

	var buf bytes.Buffer
	err := RunGroupBy(strings.NewReader(input), &buf, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "debug") {
		t.Errorf("expected 'debug' in output, got: %s", out)
	}
}
