package filter

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeLookupFile(t *testing.T, records []map[string]interface{}) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "lookup-*.ndjson")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, r := range records {
		if err := enc.Encode(r); err != nil {
			t.Fatalf("encode: %v", err)
		}
	}
	return filepath.Clean(f.Name())
}

func TestNewJoiner_Valid(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "1", "name": "alice"},
	})
	_, err := NewJoiner("id", path, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewJoiner_EmptyKey(t *testing.T) {
	_, err := NewJoiner("", "somefile", false)
	if err == nil {
		t.Fatal("expected error for empty keyField")
	}
}

func TestNewJoiner_MissingFile(t *testing.T) {
	_, err := NewJoiner("id", "/nonexistent/file.ndjson", false)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestJoiner_Apply_Match(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "42", "region": "us-east"},
	})
	j, _ := NewJoiner("id", path, false)
	rec := map[string]interface{}{"id": "42", "level": "info"}
	out := j.Apply(rec)
	if out["region"] != "us-east" {
		t.Errorf("expected region us-east, got %v", out["region"])
	}
	if out["level"] != "info" {
		t.Errorf("original field should be preserved")
	}
}

func TestJoiner_Apply_NoMatch(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "99", "region": "eu-west"},
	})
	j, _ := NewJoiner("id", path, false)
	rec := map[string]interface{}{"id": "1", "level": "warn"}
	out := j.Apply(rec)
	if _, ok := out["region"]; ok {
		t.Error("region should not be present when no match")
	}
}

func TestJoiner_Apply_NoOverwrite(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "7", "level": "debug"},
	})
	j, _ := NewJoiner("id", path, false)
	rec := map[string]interface{}{"id": "7", "level": "error"}
	out := j.Apply(rec)
	if out["level"] != "error" {
		t.Errorf("overwrite=false: expected original level 'error', got %v", out["level"])
	}
}

func TestJoiner_Apply_WithOverwrite(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"id": "7", "level": "debug"},
	})
	j, _ := NewJoiner("id", path, true)
	rec := map[string]interface{}{"id": "7", "level": "error"}
	out := j.Apply(rec)
	if out["level"] != "debug" {
		t.Errorf("overwrite=true: expected lookup level 'debug', got %v", out["level"])
	}
}

func TestRunJoinPipeline(t *testing.T) {
	path := writeLookupFile(t, []map[string]interface{}{
		{"user_id": "u1", "email": "a@example.com"},
	})
	j, _ := NewJoiner("user_id", path, false)

	input := `{"user_id":"u1","msg":"hello"}
{"user_id":"u2","msg":"world"}
not-json
`
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)
	n, err := RunJoinPipeline(bytes.NewBufferString(input), j, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 written, got %d", n)
	}
}
