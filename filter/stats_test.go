package filter

import (
	"testing"
)

func TestStats_Add(t *testing.T) {
	var s Stats
	s.Add(true, false)
	s.Add(false, false)
	s.Add(false, true)

	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Matched != 1 {
		t.Errorf("expected Matched=1, got %d", s.Matched)
	}
	if s.Skipped != 1 {
		t.Errorf("expected Skipped=1, got %d", s.Skipped)
	}
	if s.Invalid != 1 {
		t.Errorf("expected Invalid=1, got %d", s.Invalid)
	}
}

func TestStats_Summary(t *testing.T) {
	s := Stats{Total: 10, Matched: 7, Skipped: 2, Invalid: 1}
	got := s.Summary()
	want := "total=10 matched=7 skipped=2 invalid=1"
	if got != want {
		t.Errorf("Summary() = %q, want %q", got, want)
	}
}

func TestStats_MatchRate(t *testing.T) {
	s := Stats{Total: 10, Matched: 6, Skipped: 3, Invalid: 1}
	rate := s.MatchRate()
	// valid = 9, matched = 6 => 0.666...
	if rate < 0.66 || rate > 0.67 {
		t.Errorf("MatchRate() = %f, want ~0.666", rate)
	}
}

func TestStats_MatchRate_NoValid(t *testing.T) {
	var s Stats
	if s.MatchRate() != 0 {
		t.Error("expected MatchRate=0 when no valid lines")
	}
}
