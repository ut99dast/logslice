package filter

import (
	"fmt"
	"strings"
)

// Selector retains only the specified fields from each log record,
// dropping all other keys. This is useful for reducing noise when
// only a subset of fields is relevant to the current investigation.
type Selector struct {
	fields map[string]struct{}
}

// NewSelector creates a Selector that keeps only the named fields.
// At least one field name must be provided, and no field name may be
// empty. Returns an error if validation fails.
func NewSelector(fields []string) (*Selector, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("selector: at least one field name is required")
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		f = strings.TrimSpace(f)
		if f == "" {
			return nil, fmt.Errorf("selector: field name must not be empty")
		}
		set[f] = struct{}{}
	}
	return &Selector{fields: set}, nil
}

// Apply returns a new record containing only the fields listed in the
// Selector. Fields that are present in the selector but absent from
// the record are silently omitted (no error is returned).
func (s *Selector) Apply(record map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(s.fields))
	for k, v := range record {
		if _, ok := s.fields[k]; ok {
			out[k] = v
		}
	}
	return out
}

// RunSelect reads newline-delimited JSON from src, applies the selector
// to each valid record, and writes the projected records to dst.
// Lines that cannot be parsed as JSON are counted as invalid and skipped.
// The function returns the number of records written and any terminal error.
func RunSelect(src []byte, fields []string) ([]map[string]interface{}, error) {
	sel, err := NewSelector(fields)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimRight(string(src), "\n"), "\n")
	var results []map[string]interface{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		record, err := ParseRecord(line)
		if err != nil {
			// skip unparseable lines, consistent with other pipeline helpers
			continue
		}
		results = append(results, sel.Apply(record))
	}

	return results, nil
}
