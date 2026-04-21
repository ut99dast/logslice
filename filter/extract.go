package filter

import (
	"fmt"
	"io"
	"strings"
)

// Extractor selects a subset of fields from a log record, outputting only those keys.
type Extractor struct {
	fields []string
}

// NewExtractor creates an Extractor that retains only the specified fields.
// fields is a comma-separated list of field names, e.g. "level,msg,ts".
func NewExtractor(fields string) (*Extractor, error) {
	if fields == "" {
		return nil, fmt.Errorf("extract: fields must not be empty")
	}
	parts := strings.Split(fields, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		result = append(result, p)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("extract: no valid field names provided")
	}
	return &Extractor{fields: result}, nil
}

// Apply returns a new record containing only the configured fields.
// Fields not present in the record are omitted silently.
func (e *Extractor) Apply(record map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(e.fields))
	for _, f := range e.fields {
		if v, ok := record[f]; ok {
			out[f] = v
		}
	}
	return out
}

// RunExtract reads JSON lines from r, retains only the specified fields, and
// writes the results to w using the provided Writer.
func RunExtract(r io.Reader, w *Writer, fields string) error {
	ext, err := NewExtractor(fields)
	if err != nil {
		return err
	}
	scanner := NewScanner(r, nil)
	for scanner.Scan() {
		record := scanner.Record()
		out := ext.Apply(record)
		if err := w.Write(out); err != nil {
			return err
		}
	}
	return scanner.Err()
}
