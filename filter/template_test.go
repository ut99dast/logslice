package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewTemplateRenderer_Valid(t *testing.T) {
	_, err := NewTemplateRenderer("hello {{.name}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewTemplateRenderer_Empty(t *testing.T) {
	_, err := NewTemplateRenderer("")
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestNewTemplateRenderer_UnclosedBrace(t *testing.T) {
	_, err := NewTemplateRenderer("hello {{.name")
	if err == nil {
		t.Fatal("expected error for unclosed brace")
	}
}

func TestNewTemplateRenderer_EmptyField(t *testing.T) {
	_, err := NewTemplateRenderer("hello {{.}}")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestTemplateRenderer_Apply_Basic(t *testing.T) {
	r, _ := NewTemplateRenderer("user={{.user}} level={{.level}}")
	rec := map[string]interface{}{"user": "alice", "level": "info"}
	got := r.Apply(rec)
	want := "user=alice level=info"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTemplateRenderer_Apply_MissingField(t *testing.T) {
	r, _ := NewTemplateRenderer("msg={{.msg}}")
	rec := map[string]interface{}{}
	got := r.Apply(rec)
	if !strings.Contains(got, "<nil>") {
		t.Errorf("expected <nil> for missing field, got %q", got)
	}
}

func TestTemplateRenderer_Apply_LiteralOnly(t *testing.T) {
	r, _ := NewTemplateRenderer("no placeholders here")
	rec := map[string]interface{}{"x": 1}
	got := r.Apply(rec)
	if got != "no placeholders here" {
		t.Errorf("got %q", got)
	}
}

func TestRunTemplate_Basic(t *testing.T) {
	input := `{"user":"bob","action":"login"}` + "\n" +
		`{"user":"alice","action":"logout"}` + "\n"
	var out bytes.Buffer
	err := RunTemplate(strings.NewReader(input), &out, "{{.user}} did {{.action}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "bob did login" {
		t.Errorf("line 0: got %q", lines[0])
	}
	if lines[1] != "alice did logout" {
		t.Errorf("line 1: got %q", lines[1])
	}
}

func TestRunTemplate_SkipsInvalidLines(t *testing.T) {
	input := "not-json\n" + `{"user":"carol"}` + "\n"
	var out bytes.Buffer
	err := RunTemplate(strings.NewReader(input), &out, "u={{.user}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 || lines[0] != "u=carol" {
		t.Errorf("unexpected output: %v", lines)
	}
}

func TestRunTemplate_InvalidTemplate(t *testing.T) {
	err := RunTemplate(strings.NewReader(""), &bytes.Buffer{}, "")
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}
