package filter

import (
	"strings"
	"testing"
)

func TestRunSplitPipeline_Basic(t *testing.T) {
	input := `{"tags":"go,rust","level":"info"}
`
	var out strings.Builder
	args := SplitArgs{
		Field:     "tags",
		Delimiter: ",",
		OutField:  "tag",
		Format:    "json",
	}
	if err := RunSplitPipeline(strings.NewReader(input), &out, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRunSplitPipeline_MissingField(t *testing.T) {
	var out strings.Builder
	args := SplitArgs{Field: "", Delimiter: ",", Format: "json"}
	if err := RunSplitPipeline(strings.NewReader(""), &out, args); err == nil {
		t.Error("expected error for empty field")
	}
}

func TestRunSplitPipeline_MissingDelimiter(t *testing.T) {
	var out strings.Builder
	args := SplitArgs{Field: "tags", Delimiter: "", Format: "json"}
	if err := RunSplitPipeline(strings.NewReader(""), &out, args); err == nil {
		t.Error("expected error for empty delimiter")
	}
}

func TestRunSplitPipeline_SkipsInvalidLines(t *testing.T) {
	input := "not-json\n{\"tags\":\"a,b\",\"x\":\"1\"}\n"
	var out strings.Builder
	args := SplitArgs{
		Field:     "tags",
		Delimiter: ",",
		OutField:  "tag",
		Format:    "json",
	}
	if err := RunSplitPipeline(strings.NewReader(input), &out, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 output lines (from valid record), got %d", len(lines))
	}
}

func TestRunSplitPipeline_EmptyInput(t *testing.T) {
	var out strings.Builder
	args := SplitArgs{
		Field:     "tags",
		Delimiter: ",",
		Format:    "json",
	}
	if err := RunSplitPipeline(strings.NewReader(""), &out, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.String() != "" {
		t.Errorf("expected empty output, got %q", out.String())
	}
}
