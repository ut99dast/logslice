package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	cases := []struct {
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"json", FormatJSON, false},
		{"", FormatJSON, false},
		{"pretty", FormatPretty, false},
		{"csv", FormatCSV, false},
		{"xml", 0, true},
	}
	for _, c := range cases {
		got, err := ParseOutputFormat(c.input)
		if c.wantErr && err == nil {
			t.Errorf("expected error for %q", c.input)
		}
		if !c.wantErr && got != c.want {
			t.Errorf("ParseOutputFormat(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}

func TestWriter_JSON(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatJSON)
	rec := map[string]interface{}{"level": "info", "msg": "hello"}
	if err := w.Write(rec); err != nil {
		t.Fatal(err)
	}
	line := strings.TrimSpace(buf.String())
	var out map[string]interface{}
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
}

func TestWriter_Pretty(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatPretty)
	rec := map[string]interface{}{"level": "warn"}
	if err := w.Write(rec); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "\n") {
		t.Error("pretty output should contain newlines")
	}
}

func TestWriter_CSV(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, FormatCSV)
	rec := map[string]interface{}{"level": "error", "msg": "oops"}
	if err := w.Write(rec); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header+data), got %d", len(lines))
	}
}
