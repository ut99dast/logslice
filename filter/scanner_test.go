package filter

import (
	"strings"
	"testing"
)

type alwaysMatch struct{}

func (alwaysMatch) Match(_ map[string]interface{}) bool { return true }

type neverMatch struct{}

func (neverMatch) Match(_ map[string]interface{}) bool { return false }

func TestScanner_AllMatch(t *testing.T) {
	input := `{"level":"info","msg":"a"}
{"level":"warn","msg":"b"}
{"level":"error","msg":"c"}
`
	s := NewScanner(strings.NewReader(input), alwaysMatch{})
	count := 0
	for {
		_, ok := s.Next()
		if !ok {
			break
		}
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 records, got %d", count)
	}
}

func TestScanner_NoneMatch(t *testing.T) {
	input := `{"level":"info"}
{"level":"warn"}
`
	s := NewScanner(strings.NewReader(input), neverMatch{})
	if _, ok := s.Next(); ok {
		t.Error("expected no records")
	}
}

func TestScanner_SkipsInvalidLines(t *testing.T) {
	input := `not json
{"level":"info","msg":"valid"}

also not json
{"level":"warn","msg":"also valid"}
`
	s := NewScanner(strings.NewReader(input), alwaysMatch{})
	count := 0
	for {
		_, ok := s.Next()
		if !ok {
			break
		}
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 valid records, got %d", count)
	}
}

func TestScanner_NilFilter(t *testing.T) {
	input := `{"level":"debug"}
`
	s := NewScanner(strings.NewReader(input), nil)
	_, ok := s.Next()
	if !ok {
		t.Error("nil filter should match all records")
	}
}
