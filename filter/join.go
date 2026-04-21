package filter

import (
	"encoding/json"
	"fmt"
	"os"
)

// Joiner enriches records from a primary stream by looking up matching
// records from a static lookup file, keyed on a shared field.
type Joiner struct {
	keyField   string
	lookupData map[string]map[string]interface{}
	overwrite  bool
}

// NewJoiner creates a Joiner that loads lookup records from lookupPath,
// indexes them by keyField, and optionally overwrites existing fields.
func NewJoiner(keyField, lookupPath string, overwrite bool) (*Joiner, error) {
	if keyField == "" {
		return nil, fmt.Errorf("join: keyField must not be empty")
	}
	f, err := os.Open(lookupPath)
	if err != nil {
		return nil, fmt.Errorf("join: cannot open lookup file: %w", err)
	}
	defer f.Close()

	lookup := make(map[string]map[string]interface{})
	dec := json.NewDecoder(f)
	for dec.More() {
		var rec map[string]interface{}
		if err := dec.Decode(&rec); err != nil {
			continue
		}
		if kv, ok := rec[keyField]; ok {
			key := fmt.Sprintf("%v", kv)
			lookup[key] = rec
		}
	}
	return &Joiner{keyField: keyField, lookupData: lookup, overwrite: overwrite}, nil
}

// Apply merges lookup fields into rec. Returns the enriched record.
// If no match is found the original record is returned unchanged.
func (j *Joiner) Apply(rec map[string]interface{}) map[string]interface{} {
	kv, ok := rec[j.keyField]
	if !ok {
		return rec
	}
	key := fmt.Sprintf("%v", kv)
	lookupRec, found := j.lookupData[key]
	if !found {
		return rec
	}
	out := shallowCopy(rec)
	for k, v := range lookupRec {
		if k == j.keyField {
			continue
		}
		if _, exists := out[k]; exists && !j.overwrite {
			continue
		}
		out[k] = v
	}
	return out
}

// RunJoin reads NDJSON lines from src, enriches each record using j,
// and writes results to dst via w.
func RunJoin(src []string, j *Joiner, w *Writer) error {
	for _, line := range src {
		rec, err := ParseRecord(line)
		if err != nil {
			continue
		}
		enriched := j.Apply(rec)
		if err := w.Write(enriched); err != nil {
			return err
		}
	}
	return nil
}
