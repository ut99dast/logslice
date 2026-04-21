package filter

import (
	"fmt"
	"io"
)

// Wrapper wraps a string field value in a configurable prefix and suffix.
type Wrapper struct {
	field  string
	prefix string
	suffix string
}

// NewWrapper creates a Wrapper for the given field with the specified prefix and suffix.
// Returns an error if field is empty.
func NewWrapper(field, prefix, suffix string) (*Wrapper, error) {
	if field == "" {
		return nil, fmt.Errorf("wrap: field name must not be empty")
	}
	return &Wrapper{field: field, prefix: prefix, suffix: suffix}, nil
}

// Apply wraps the field value in the record. If the field is missing or not a
// string, the record is returned unchanged.
func (w *Wrapper) Apply(rec Record) Record {
	val, ok := rec[w.field]
	if !ok {
		return rec
	}
	str, ok := val.(string)
	if !ok {
		return rec
	}
	out := shallowCopy(rec)
	out[w.field] = w.prefix + str + w.suffix
	return out
}

// RunWrap reads JSON lines from r, wraps the specified field value using the
// given prefix and suffix, and writes results to w.
func RunWrap(r io.Reader, w io.Writer, field, prefix, suffix string, format OutputFormat) error {
	wrapper, err := NewWrapper(field, prefix, suffix)
	if err != nil {
		return err
	}
	writer, err := NewWriter(w, format)
	if err != nil {
		return err
	}
	scanner := NewScanner(r)
	for scanner.Scan() {
		rec, err := scanner.Record()
		if err != nil {
			continue
		}
		if err := writer.Write(wrapper.Apply(rec)); err != nil {
			return err
		}
	}
	return writer.Flush()
}
