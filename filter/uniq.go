package filter

import (
	"fmt"
	"io"
)

// Uniquer tracks consecutive duplicate values for a given field
// and emits only the first occurrence of each run.
type Uniquer struct {
	field    string
	lastVal  interface{}
	hasLast  bool
}

// NewUniqer creates a Uniquer that filters consecutive duplicates on field.
// Pass an empty string to compare whole records.
func NewUniqer(field string) (*Uniquer, error) {
	if field == "" {
		return nil, fmt.Errorf("uniq: field name must not be empty")
	}
	return &Uniquer{field: field}, nil
}

// Apply returns (record, true) if the field value differs from the previous
// record, or (nil, false) if it is a consecutive duplicate.
func (u *Uniquer) Apply(record map[string]interface{}) (map[string]interface{}, bool) {
	val, ok := record[u.field]
	if !ok {
		// field missing — always pass through
		u.hasLast = false
		return record, true
	}
	if u.hasLast && fmt.Sprintf("%v", val) == fmt.Sprintf("%v", u.lastVal) {
		return nil, false
	}
	u.lastVal = val
	u.hasLast = true
	return record, true
}

// Reset clears the remembered last value.
func (u *Uniquer) Reset() {
	u.lastVal = nil
	u.hasLast = false
}

// RunUniq reads JSON lines from r, suppresses consecutive duplicate values
// for field, and writes surviving records to w.
func RunUniq(r io.Reader, w io.Writer, field string) error {
	uniqer, err := NewUniqer(field)
	if err != nil {
		return err
	}
	scanner := NewScanner(r)
	writer := NewWriter(w, FormatJSON)
	for scanner.Scan() {
		rec, err := ParseRecord(scanner.Text())
		if err != nil {
			continue
		}
		out, keep := uniqer.Apply(rec)
		if keep {
			if err := writer.Write(out); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
