package filter

import (
	"bufio"
	"io"
)

// Scanner reads log records line by line from an io.Reader,
// parses each line as JSON, and emits matching records.
type Scanner struct {
	reader  io.Reader
	filter  Filter
	scanner *bufio.Scanner
}

// Filter is the interface implemented by all log filters.
type Filter interface {
	Match(record map[string]interface{}) bool
}

// NewScanner creates a Scanner that applies the given filter.
func NewScanner(r io.Reader, f Filter) *Scanner {
	return &Scanner{
		reader:  r,
		filter:  f,
		scanner: bufio.NewScanner(r),
	}
}

// Next advances to the next matching record.
// Returns the record and true, or nil and false when done.
func (s *Scanner) Next() (map[string]interface{}, bool) {
	for s.scanner.Scan() {
		line := s.scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		record, err := ParseRecord(line)
		if err != nil {
			continue
		}
		if s.filter == nil || s.filter.Match(record) {
			return record, true
		}
	}
	return nil, false
}

// Err returns any scanner error after Next returns false.
func (s *Scanner) Err() error {
	return s.scanner.Err()
}
