package filter

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLimiter_Invalid(t *testing.T) {
	_, err := NewLimiter(-1)
	if err == nil {
		t.Fatal("expected error for negative limit")
	}
}

func TestLimiter_Unlimited(t *testing.T) {
	l, _ := NewLimiter(0)
	for i := 0; i < 1000; i++ {
		if l.Done() {
			t.Fatal("unlimited limiter should never be done")
		}
		l.Accept()
	}
}

func TestLimiter_AcceptAndDone(t *testing.T) {
	l, _ := NewLimiter(3)
	for i := 0; i < 3; i++ {
		if l.Done() {
			t.Fatalf("should not be done at iteration %d", i)
		}
		if !l.Accept() {
			t.Fatalf("Accept should return true at iteration %d", i)
		}
	}
	if !l.Done() {
		t.Fatal("expected Done after 3 accepts")
	}
	if l.Accept() {
		t.Fatal("Accept should return false when limit reached")
	}
}

func TestLimiter_Reset(t *testing.T) {
	l, _ := NewLimiter(2)
	l.Accept()
	l.Accept()
	if !l.Done() {
		t.Fatal("expected Done")
	}
	l.Reset()
	if l.Done() {
		t.Fatal("expected not Done after Reset")
	}
}

func TestRunLimit(t *testing.T) {
	input := strings.Join([]string{
		`{"level":"info","msg":"a"}`,
		`{"level":"info","msg":"b"}`,
		`{"level":"info","msg":"c"}`,
		`{"level":"info","msg":"d"}`,
	}, "\n")

	scanner := NewScanner(strings.NewReader(input))
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)

	err := RunLimit(scanner, nil, 2, w)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestRunLimit_Zero_NoLimit(t *testing.T) {
	input := strings.Join([]string{
		`{"msg":"1"}`,
		`{"msg":"2"}`,
		`{"msg":"3"}`,
	}, "\n")
	scanner := NewScanner(strings.NewReader(input))
	var buf bytes.Buffer
	w, _ := NewWriter(&buf, FormatJSON)

	_ = RunLimit(scanner, nil, 0, w)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}
