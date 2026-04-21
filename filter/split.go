package filter

import (
	"fmt"
	"strings"
)

// Splitter splits a string field into multiple records, one per token.
type Splitter struct {
	field     string
	delimiter string
	outField  string
}

// NewSplitter creates a Splitter that splits field by delimiter.
// outField is the name of the field in each output record holding the token.
// If outField is empty, it defaults to field.
func NewSplitter(field, delimiter, outField string) (*Splitter, error) {
	if field == "" {
		return nil, fmt.Errorf("split: field name must not be empty")
	}
	if delimiter == "" {
		return nil, fmt.Errorf("split: delimiter must not be empty")
	}
	if outField == "" {
		outField = field
	}
	return &Splitter{field: field, delimiter: delimiter, outField: outField}, nil
}

// Apply returns one record per token found by splitting record[field] by delimiter.
// The original field is replaced by the token in each output record.
// Returns an error if the field is missing or not a string.
func (s *Splitter) Apply(record map[string]interface{}) ([]map[string]interface{}, error) {
	v, ok := record[s.field]
	if !ok {
		return nil, fmt.Errorf("split: field %q not found", s.field)
	}
	str, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("split: field %q is not a string", s.field)
	}
	tokens := strings.Split(str, s.delimiter)
	results := make([]map[string]interface{}, 0, len(tokens))
	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		out := shallowCopy(record)
		delete(out, s.field)
		out[s.outField] = tok
		results = append(results, out)
	}
	return results, nil
}

// RunSplit reads JSON lines from scanner, splits each record on field/delimiter,
// and writes resulting records to w.
func RunSplit(scanner *Scanner, w *Writer, field, delimiter, outField string) error {
	splitter, err := NewSplitter(field, delimiter, outField)
	if err != nil {
		return err
	}
	for scanner.Scan() {
		record, err := scanner.Record()
		if err != nil {
			continue
		}
		records, err := splitter.Apply(record)
		if err != nil {
			continue
		}
		for _, r := range records {
			if werr := w.Write(r); werr != nil {
				return werr
			}
		}
	}
	return scanner.Err()
}
