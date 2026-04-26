package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewWindower_Invalid(t *testing.T) {
	_, err := NewWindower("", 0, 0)
	if err == nil {
		t.Fatal("expected error for size=0")
	}
	_, err = NewWindower("", 3, 3)
	if err == nil {
		t.Fatal("expected error for overlap >= size")
	}
	_, err = NewWindower("", 3, -1)
	if err == nil {
		t.Fatal("expected error for negative overlap")
	}
}

func TestWindower_Tumbling(t *testing.T) {
	w, _ := NewWindower("", 2, 0)
	records := []map[string]interface{}{
		{"n": "1"}, {"n": "2"}, {"n": "3"}, {"n": "4"}, {"n": "5"},
	}
	groups := w.Apply(records)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if len(groups[0].Records) != 2 {
		t.Errorf("group 0 should have 2 records")
	}
	if len(groups[2].Records) != 1 {
		t.Errorf("group 2 (last) should have 1 record")
	}
}

func TestWindower_Sliding(t *testing.T) {
	w, _ := NewWindower("", 3, 1)
	records := []map[string]interface{}{
		{"n": "a"}, {"n": "b"}, {"n": "c"}, {"n": "d"},
	}
	groups := w.Apply(records)
	// step = 3-1 = 2 → windows at [0,3), [2,4)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[1].Records[0]["n"] != "c" {
		t.Errorf("second window should start at 'c'")
	}
}

func TestWindower_Empty(t *testing.T) {
	w, _ := NewWindower("", 3, 0)
	groups := w.Apply(nil)
	if len(groups) != 0 {
		t.Errorf("expected no groups for empty input")
	}
}

func TestRunWindow_Basic(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"a"}`,
		`{"level":"warn","msg":"b"}`,
		`{"level":"error","msg":"c"}`,
		`not json`,
	}, "\n")

	var out bytes.Buffer
	err := RunWindow(strings.NewReader(input), &out, "", 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 window lines, got %d", len(lines))
	}

	var first map[string]interface{}
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("invalid JSON on first window: %v", err)
	}
	if first["window"] != "0" {
		t.Errorf("expected window index 0")
	}
}
