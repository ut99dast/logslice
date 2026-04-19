package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestAggregateResult_Add(t *testing.T) {
	a := NewAggregateResult("level")
	a.Add(map[string]interface{}{"level": "info"})
	a.Add(map[string]interface{}{"level": "info"})
	a.Add(map[string]interface{}{"level": "error"})
	a.Add(map[string]interface{}{"msg": "no level"})

	if a.Counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", a.Counts["info"])
	}
	if a.Counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", a.Counts["error"])
	}
	if a.Counts["<missing>"] != 1 {
		t.Errorf("expected <missing>=1, got %d", a.Counts["<missing>"])
	}
}

func TestAggregateResult_Sorted(t *testing.T) {
	a := NewAggregateResult("level")
	a.Counts["debug"] = 1
	a.Counts["info"] = 5
	a.Counts["error"] = 3

	entries := a.Sorted()
	if entries[0].Value != "info" {
		t.Errorf("expected first entry info, got %s", entries[0].Value)
	}
	if entries[1].Value != "error" {
		t.Errorf("expected second entry error, got %s", entries[1].Value)
	}
}

func TestRunAggregate(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"started"}`,
		`{"level":"info","msg":"running"}`,
		`{"level":"error","msg":"failed"}`,
		`not json`,
	}, "\n")

	var buf bytes.Buffer
	err := RunAggregate(strings.NewReader(input), &buf, "level", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "info") || !strings.Contains(out, "2") {
		t.Errorf("expected info count 2 in output:\n%s", out)
	}
	if !strings.Contains(out, "error") {
		t.Errorf("expected error in output:\n%s", out)
	}
}
