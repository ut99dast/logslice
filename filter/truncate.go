package filter

import (
	"fmt"
	"io"
)

// Truncator truncates string field values to a maximum length.
type Truncator struct {
	field  string
	maxLen int
	suffix string
}

// NewTruncator creates a Truncator for the given field and max length.
// suffix is appended when a value is truncated (e.g. "...").
func NewTruncator(field string, maxLen int, suffix string) (*Truncator, error) {
	if field == "" {
		return nil, fmt.Errorf("truncate: field name must not be empty")
	}
	if maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be positive, got %d", maxLen)
	}
	return &Truncator{field: field, maxLen: maxLen, suffix: suffix}, nil
}

// Apply truncates the configured field in the record if its string value
// exceeds maxLen. Records missing the field are passed through unchanged.
func (t *Truncator) Apply(rec Record) (Record, error) {
	out := make(Record, len(rec))
	for k, v := range rec {
		out[k] = v
	}
	val, ok := out[t.field]
	if !ok {
		return out, nil
	}
	s, ok := val.(string)
	if !ok {
		return out, nil
	}
	if len(s) > t.maxLen {
		out[t.field] = s[:t.maxLen] + t.suffix
	}
	return out, nil
}

// RunTruncate reads JSON records from r, truncates the specified field,
// and writes results to w. Invalid lines are skipped and counted.
func RunTruncate(r io.Reader, w io.Writer, field string, maxLen int, suffix string, format OutputFormat) error {
	tr, err := NewTruncator(field, maxLen, suffix)
	if err != nil {
		return err
	}
	writer, err := NewWriter(w, format)
	if err != nil {
		return err
	}
	scanner := NewScanner(r, MatchAll())
	for scanner.Scan() {
		rec := scanner.Record()
		result, err := tr.Apply(rec)
		if err != nil {
			continue
		}
		if err := writer.Write(result); err != nil {
			return err
		}
	}
	return writer.Flush()
}
