package filter

import (
	"fmt"
	"io"
)

// Coalescer returns the first non-empty value from a list of fields,
// writing the result into an output field.
type Coalescer struct {
	fields    []string
	outField  string
}

// NewCoalescer creates a Coalescer that checks each field in order and
// writes the first non-empty value to outField.
// At least two source fields are required.
func NewCoalescer(fields []string, outField string) (*Coalescer, error) {
	if len(fields) < 2 {
		return nil, fmt.Errorf("coalesce requires at least 2 source fields")
	}
	if outField == "" {
		return nil, fmt.Errorf("coalesce output field must not be empty")
	}
	return &Coalescer{fields: fields, outField: outField}, nil
}

// Apply scans fields in order and sets outField to the first non-empty
// string value found. If no field has a non-empty value, outField is set
// to an empty string. The original record is not mutated.
func (c *Coalescer) Apply(rec map[string]interface{}) map[string]interface{} {
	out := shallowCopy(rec)
	for _, f := range c.fields {
		val, ok := rec[f]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if ok && s != "" {
			out[c.outField] = s
			return out
		}
		// non-string non-nil values also count
		if !ok && val != nil {
			out[c.outField] = val
			return out
		}
	}
	out[c.outField] = ""
	return out
}

// RunCoalesce reads newline-delimited JSON from r, applies the coalescer
// to each valid record, and writes results to w.
func RunCoalesce(r io.Reader, w io.Writer, fields []string, outField string, format string) error {
	c, err := NewCoalescer(fields, outField)
	if err != nil {
		return err
	}
	writer, err := NewWriter(w, format)
	if err != nil {
		return err
	}
	scanner := NewScanner(r)
	for scanner.Scan() {
		rec, err := ParseRecord(scanner.Bytes())
		if err != nil {
			continue
		}
		if err := writer.Write(c.Apply(rec)); err != nil {
			return err
		}
	}
	return writer.Flush()
}
