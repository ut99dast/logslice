package filter

import (
	"encoding/json"
	"fmt"
	"io"
)

// Merger merges fields from a static JSON object into each log record.
// Fields from the overlay are added to every record; existing fields
// are overwritten only when overwrite is true.
type Merger struct {
	overlay   map[string]interface{}
	overwrite bool
}

// NewMerger creates a Merger from a JSON string representing the fields to
// merge into every record. Set overwrite=true to allow overlay fields to
// replace existing record fields.
func NewMerger(jsonFields string, overwrite bool) (*Merger, error) {
	if jsonFields == "" {
		return nil, fmt.Errorf("merge: fields JSON must not be empty")
	}
	var overlay map[string]interface{}
	if err := json.Unmarshal([]byte(jsonFields), &overlay); err != nil {
		return nil, fmt.Errorf("merge: invalid JSON fields: %w", err)
	}
	if len(overlay) == 0 {
		return nil, fmt.Errorf("merge: fields JSON must contain at least one field")
	}
	return &Merger{overlay: overlay, overwrite: overwrite}, nil
}

// Apply merges the overlay fields into rec and returns the updated record.
func (m *Merger) Apply(rec map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(rec)+len(m.overlay))
	for k, v := range rec {
		out[k] = v
	}
	for k, v := range m.overlay {
		if _, exists := out[k]; !exists || m.overwrite {
			out[k] = v
		}
	}
	return out
}

// RunMerge reads JSON lines from r, merges overlay fields into each valid
// record, and writes the result to w. Invalid lines are skipped.
func RunMerge(r io.Reader, w io.Writer, jsonFields string, overwrite bool) error {
	merger, err := NewMerger(jsonFields, overwrite)
	if err != nil {
		return err
	}
	scanner := NewScanner(r)
	encoder := json.NewEncoder(w)
	for scanner.Scan() {
		rec, err := ParseRecord(scanner.Bytes())
		if err != nil {
			continue
		}
		merged := merger.Apply(rec)
		if err := encoder.Encode(merged); err != nil {
			return fmt.Errorf("merge: encode error: %w", err)
		}
	}
	return scanner.Err()
}
