package filter

import "time"

// TimeRange holds an optional start and end boundary.
type TimeRange struct {
	Start *time.Time
	End   *time.Time
}

// NewTimeRange constructs a TimeRange from optional RFC3339 strings.
func NewTimeRange(start, end string) (TimeRange, error) {
	var tr TimeRange
	if start != "" {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return tr, fmt.Errorf("invalid start time: %w", err)
		}
		tr.Start = &t
	}
	if end != "" {
		t, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return tr, fmt.Errorf("invalid end time: %w", err)
		}
		tr.End = &t
	}
	return tr, nil
}

// Match returns true when t falls within the TimeRange.
func (tr TimeRange) Match(t time.Time) bool {
	if tr.Start != nil && t.Before(*tr.Start) {
		return false
	}
	if tr.End != nil && t.After(*tr.End) {
		return false
	}
	return true
}
