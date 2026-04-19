package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewSorter_Invalid(t *testing.T) {
	_, err := NewSorter("", "asc")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
	_, err = NewSorter("level", "sideways")
	if err == nil {
		t.Fatal("expected error for unknown order")
	}
}

func TestSorter_Asc(t *testing.T) {
	s, _ := NewSorter("level", "asc")
	s.Add(map[string]interface{}{"level": "warn", "msg": "b"})
	s.Add(map[string]interface{}{"level": "error", "msg": "c"})
	s.Add(map[string]interface{}{"level": "info", "msg": "a"})

	var buf bytes.Buffer
	if err := s.Sort(&buf); err != nil {
		t.Fatalf("Sort error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	var first map[string]interface{}
	json.Unmarshal([]byte(lines[0]), &first)
	if first["level"] != "error" {
		t.Errorf("expected first level=error, got %v", first["level"])
	}
}

func TestSorter_Desc(t *testing.T) {
	s, _ := NewSorter("level", "desc")
	s.Add(map[string]interface{}{"level": "info"})
	s.Add(map[string]interface{}{"level": "warn"})

	var buf bytes.Buffer
	s.Sort(&buf)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	var first map[string]interface{}
	json.Unmarshal([]byte(lines[0]), &first)
	if first["level"] != "warn" {
		t.Errorf("expected first level=warn, got %v", first["level"])
	}
}

func TestRunSort(t *testing.T) {
	input := `{"ts":"2024-01-03","msg":"c"}
{"ts":"2024-01-01","msg":"a"}
{"ts":"2024-01-02","msg":"b"}
`
	var buf bytes.Buffer
	if err := RunSort(strings.NewReader(input), &buf, "ts", "asc"); err != nil {
		t.Fatalf("RunSort error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	var first map[string]interface{}
	json.Unmarshal([]byte(lines[0]), &first)
	if first["msg"] != "a" {
		t.Errorf("expected first msg=a, got %v", first["msg"])
	}
}
