package filter

import (
	"strings"
	"testing"
)

func TestNewHeader_Invalid(t *testing.T) {
	_, err := NewHeader(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
	_, err = NewHeader(-5)
	if err == nil {
		t.Fatal("expected error for n=-5")
	}
}

func TestHeader_FewerThanN(t *testing.T) {
	input := `{"level":"info","msg":"a"}
{"level":"info","msg":"b"}
`
	h, _ := NewHeader(5)
	var out strings.Builder
	count, err := h.Take(strings.NewReader(input), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 records, got %d", count)
	}
}

func TestHeader_ExactlyN(t *testing.T) {
	input := `{"msg":"1"}
{"msg":"2"}
{"msg":"3"}
`
	h, _ := NewHeader(3)
	var out strings.Builder
	count, _ := h.Take(strings.NewReader(input), &out)
	if count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestHeader_MoreThanN(t *testing.T) {
	input := `{"msg":"1"}
{"msg":"2"}
{"msg":"3"}
{"msg":"4"}
`
	h, _ := NewHeader(2)
	var out strings.Builder
	count, _ := h.Take(strings.NewReader(input), &out)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestHeader_SkipsInvalidLines(t *testing.T) {
	input := "not json\n{\"msg\":\"ok\"}\nalso not json\n{\"msg\":\"ok2\"}\n"
	h, _ := NewHeader(2)
	var out strings.Builder
	count, _ := h.Take(strings.NewReader(input), &out)
	if count != 2 {
		t.Fatalf("expected 2 valid records, got %d", count)
	}
}

func TestRunHead(t *testing.T) {
	input := `{"msg":"1"}
{"msg":"2"}
{"msg":"3"}
`
	var out strings.Builder
	count, err := RunHead(strings.NewReader(input), &out, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}
