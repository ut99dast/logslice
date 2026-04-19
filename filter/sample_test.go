package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewSampler_InvalidRate(t *testing.T) {
	_, err := NewSampler(0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestSampler_KeepEveryOther(t *testing.T) {
	s, _ := NewSampler(2)
	results := make([]bool, 6)
	for i := range results {
		results[i] = s.Keep()
	}
	// positions 0,2,4 should be true
	for i, v := range results {
		expected := i%2 == 0
		if v != expected {
			t.Errorf("index %d: got %v want %v", i, v, expected)
		}
	}
}

func TestSampler_Reset(t *testing.T) {
	s, _ := NewSampler(3)
	s.Keep()
	s.Keep()
	s.Reset()
	if !s.Keep() {
		t.Error("expected first keep after reset to be true")
	}
}

func TestRunSample_RateOne(t *testing.T) {
	input := `{"level":"info","msg":"a"}
{"level":"info","msg":"b"}
{"level":"info","msg":"c"}
`
	var buf bytes.Buffer
	if err := RunSample(strings.NewReader(input), &buf, SampleConfig{Rate: 1}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestRunSample_RateTwo(t *testing.T) {
	input := `{"level":"info","msg":"a"}
{"level":"info","msg":"b"}
{"level":"info","msg":"c"}
{"level":"info","msg":"d"}
`
	var buf bytes.Buffer
	if err := RunSample(strings.NewReader(input), &buf, SampleConfig{Rate: 2}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestRunSample_InvalidLines(t *testing.T) {
	input := "not-json\n{\"msg\":\"ok\"}\n"
	var buf bytes.Buffer
	if err := RunSample(strings.NewReader(input), &buf, SampleConfig{Rate: 1}); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}
