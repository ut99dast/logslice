package filter

import "fmt"

// Stats tracks processing statistics for a log slicing run.
type Stats struct {
	Total   int
	Matched int
	Skipped int
	Invalid int
}

// Add records a line result into the stats.
func (s *Stats) Add(matched, invalid bool) {
	s.Total++
	if invalid {
		s.Invalid++
		return
	}
	if matched {
		s.Matched++
	} else {
		s.Skipped++
	}
}

// Summary returns a human-readable summary string.
func (s *Stats) Summary() string {
	return fmt.Sprintf(
		"total=%d matched=%d skipped=%d invalid=%d",
		s.Total, s.Matched, s.Skipped, s.Invalid,
	)
}

// MatchRate returns the fraction of valid lines that matched, 0 if none.
func (s *Stats) MatchRate() float64 {
	valid := s.Total - s.Invalid
	if valid == 0 {
		return 0
	}
	return float64(s.Matched) / float64(valid)
}
