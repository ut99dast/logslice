package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewCounter_Valid(t *testing.T) {
	c, err := NewCounter("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil counter")
	}
}

func TestNewCounter_EmptyField(t *testing.T) {
	_, err := NewCounter("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestCounter_Add(t *testing.T) {
	c, _ := NewCounter("level")

	records := []map[string]interface{}{
		{"level": "info"},
		{"level": "error"},
		{"level": "info"},
		{"msg": "no level field"},
	}
	for _, r := range records {
		c.Add(r)
	}

	results := c.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
	if results[0].Value != "info" || results[0].Count != 2 {
		t.Errorf("expected info=2, got %s=%d", results[0].Value, results[0].Count)
	}
	if results[1].Value != "error" || results[1].Count != 1 {
		t.Errorf("expected error=1, got %s=%d", results[1].Value, results[1].Count)
	}
}

func TestCounter_SortedDescending(t *testing.T) {
	c, _ := NewCounter("status")
	for i := 0; i < 5; i++ {
		c.Add(map[string]interface{}{"status": "ok"})
	}
	for i := 0; i < 2; i++ {
		c.Add(map[string]interface{}{"status": "fail"})
	}
	c.Add(map[string]interface{}{"status": "timeout"})

	results := c.Results()
	if results[0].Value != "ok" {
		t.Errorf("expected ok first, got %s", results[0].Value)
	}
	if results[len(results)-1].Count > results[0].Count {
		t.Error("results not sorted descending")
	}
}

func TestRunCount(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"started"}`,
		`{"level":"error","msg":"failed"}`,
		`{"level":"info","msg":"done"}`,
		`not valid json`,
	}, "\n")

	var buf bytes.Buffer
	err := RunCount(strings.NewReader(input), &buf, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}

	var first map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("invalid JSON in output: %v", err)
	}
	if first["level"] != "info" {
		t.Errorf("expected first entry to be info, got %v", first["level"])
	}
	if int(first["count"].(float64)) != 2 {
		t.Errorf("expected count 2, got %v", first["count"])
	}
}

func TestRunCount_EmptyField(t *testing.T) {
	err := RunCount(strings.NewReader(""), &bytes.Buffer{}, "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
