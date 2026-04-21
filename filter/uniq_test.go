package filter

import (
	"strings"
	"testing"
)

func TestNewUniqer_Valid(t *testing.T) {
	_, err := NewUniqer("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewUniqer_EmptyField(t *testing.T) {
	_, err := NewUniqer("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestUniqer_NoDuplicates(t *testing.T) {
	u, _ := NewUniqer("level")
	records := []map[string]interface{}{
		{"level": "info"},
		{"level": "warn"},
		{"level": "error"},
	}
	var kept int
	for _, rec := range records {
		_, ok := u.Apply(rec)
		if ok {
			kept++
		}
	}
	if kept != 3 {
		t.Errorf("expected 3 kept, got %d", kept)
	}
}

func TestUniqer_ConsecutiveDuplicates(t *testing.T) {
	u, _ := NewUniqer("level")
	records := []map[string]interface{}{
		{"level": "info"},
		{"level": "info"},
		{"level": "warn"},
		{"level": "warn"},
		{"level": "warn"},
		{"level": "info"},
	}
	var kept int
	for _, rec := range records {
		_, ok := u.Apply(rec)
		if ok {
			kept++
		}
	}
	// expect: info, warn, info => 3
	if kept != 3 {
		t.Errorf("expected 3 kept, got %d", kept)
	}
}

func TestUniqer_MissingField(t *testing.T) {
	u, _ := NewUniqer("level")
	rec := map[string]interface{}{"msg": "hello"}
	_, ok := u.Apply(rec)
	if !ok {
		t.Error("expected record with missing field to pass through")
	}
}

func TestUniqer_Reset(t *testing.T) {
	u, _ := NewUniqer("level")
	u.Apply(map[string]interface{}{"level": "info"})
	u.Reset()
	_, ok := u.Apply(map[string]interface{}{"level": "info"})
	if !ok {
		t.Error("expected record to pass through after reset")
	}
}

func TestRunUniq_Basic(t *testing.T) {
	input := `{"level":"info","msg":"a"}
{"level":"info","msg":"b"}
{"level":"warn","msg":"c"}
{"level":"warn","msg":"d"}
{"level":"info","msg":"e"}
`
	var out strings.Builder
	err := RunUniq(strings.NewReader(input), &out, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 output lines, got %d: %v", len(lines), lines)
	}
}

func TestRunUniq_InvalidField(t *testing.T) {
	err := RunUniq(strings.NewReader(""), &strings.Builder{}, "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
