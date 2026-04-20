package filter

import (
	"strings"
	"testing"
)

func TestNewMatcher_Valid(t *testing.T) {
	m, err := NewMatcher("msg", `error.*timeout`, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Matcher")
	}
}

func TestNewMatcher_Invalid(t *testing.T) {
	cases := []struct {
		field, pattern string
	}{
		{"", "foo"},
		{"msg", ""},
		{"msg", "[invalid"},
	}
	for _, tc := range cases {
		_, err := NewMatcher(tc.field, tc.pattern, false)
		if err == nil {
			t.Errorf("expected error for field=%q pattern=%q", tc.field, tc.pattern)
		}
	}
}

func TestMatcher_Match(t *testing.T) {
	m, _ := NewMatcher("msg", `^error`, false)

	cases := []struct {
		record map[string]interface{}
		want   bool
	}{
		{map[string]interface{}{"msg": "error: disk full"}, true},
		{map[string]interface{}{"msg": "info: all good"}, false},
		{map[string]interface{}{"level": "error"}, false}, // field missing
	}
	for _, tc := range cases {
		if got := m.Match(tc.record); got != tc.want {
			t.Errorf("Match(%v) = %v, want %v", tc.record, got, tc.want)
		}
	}
}

func TestMatcher_Invert(t *testing.T) {
	m, _ := NewMatcher("level", `^debug$`, true)

	if m.Match(map[string]interface{}{"level": "debug"}) {
		t.Error("inverted matcher should not match 'debug'")
	}
	if !m.Match(map[string]interface{}{"level": "error"}) {
		t.Error("inverted matcher should match non-debug level")
	}
	// missing field with invert should return true
	if !m.Match(map[string]interface{}{"msg": "hello"}) {
		t.Error("inverted matcher with missing field should return true")
	}
}

func TestRunRegex(t *testing.T) {
	input := strings.Join([]string{
		`{"msg":"error: timeout","level":"error"}`,
		`{"msg":"info: started","level":"info"}`,
		`{"msg":"error: disk full","level":"error"}`,
		`not valid json`,
	}, "\n")

	scanner := NewScanner(strings.NewReader(input))
	var out strings.Builder
	writer, _ := NewWriter(&out, FormatJSON)
	matcher, _ := NewMatcher("msg", `^error`, false)

	total, matched, err := RunRegex(scanner, writer, matcher)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 4 {
		t.Errorf("total = %d, want 4", total)
	}
	if matched != 2 {
		t.Errorf("matched = %d, want 2", matched)
	}
}
