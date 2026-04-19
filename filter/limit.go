package filter

import (
	"errors"
	"io"
)

// Limiter tracks how many records have been emitted and signals
// when the configured maximum has been reached.
type Limiter struct {
	max  int
	seen int
}

// NewLimiter creates a Limiter that allows at most n records through.
// n <= 0 is treated as unlimited (always returns false for Done).
func NewLimiter(n int) (*Limiter, error) {
	if n < 0 {
		return nil, errors.New("limit must be >= 0")
	}
	return &Limiter{max: n}, nil
}

// Done reports whether the limit has been reached.
func (l *Limiter) Done() bool {
	if l.max == 0 {
		return false
	}
	return l.seen >= l.max
}

// Accept increments the seen counter and returns true if the record
// should be kept (i.e. the limit has not yet been exceeded).
func (l *Limiter) Accept() bool {
	if l.Done() {
		return false
	}
	l.seen++
	return true
}

// Reset clears the seen counter.
func (l *Limiter) Reset() {
	l.seen = 0
}

// RunLimit reads JSON lines from scanner, applies the filter chain, and
// writes at most limit records to w. Pass limit=0 for no cap.
func RunLimit(scanner *Scanner, filters []Filter, limit int, w *Writer) error {
	lim, err := NewLimiter(limit)
	if err != nil {
		return err
	}
	for {
		record, err := scanner.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if lim.Done() {
			break
		}
		passed := true
		for _, f := range filters {
			if !f.Match(record) {
				passed = false
				break
			}
		}
		if !passed {
			continue
		}
		if !lim.Accept() {
			break
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	return nil
}
