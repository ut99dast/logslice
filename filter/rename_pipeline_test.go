package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunRenamePipeline_Basic(t *testing.T) {
	input := `{"level":"info","msg":"ok"}
{"level":"warn","msg":"slow"}
`
	var out bytes.Buffer
	cfg := RenameConfig{
		Mappings: []string{"level=severity", "msg=message"},
		Format:   "json",
	}
	if err := RunRenamePipeline(strings.NewReader(input), &out, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if strings.Contains(result, "\"level\"") {
		t.Error("output should not contain old field 'level'")
	}
	if !strings.Contains(result, "\"severity\"") {
		t.Error("output should contain new field 'severity'")
	}
	if strings.Contains(result, "\"msg\"") {
		t.Error("output should not contain old field 'msg'")
	}
	if !strings.Contains(result, "\"message\"") {
		t.Error("output should contain new field 'message'")
	}
}

func TestRunRenamePipeline_NoMappings(t *testing.T) {
	var out bytes.Buffer
	cfg := RenameConfig{Mappings: nil, Format: "json"}
	err := RunRenamePipeline(strings.NewReader("{}"), &out, cfg)
	if err == nil {
		t.Error("expected error for empty mappings")
	}
}

func TestRunRenamePipeline_InvalidFormat(t *testing.T) {
	var out bytes.Buffer
	cfg := RenameConfig{
		Mappings: []string{"a=b"},
		Format:   "xml",
	}
	err := RunRenamePipeline(strings.NewReader("{}"), &out, cfg)
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestRunRenamePipeline_SkipsInvalidLines(t *testing.T) {
	input := "not json\n{\"level\":\"info\"}\nalso not json\n"
	var out bytes.Buffer
	cfg := RenameConfig{
		Mappings: []string{"level=severity"},
		Format:   "json",
	}
	if err := RunRenamePipeline(strings.NewReader(input), &out, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 output line, got %d", len(lines))
	}
}
